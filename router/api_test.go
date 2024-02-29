package router_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"

	"file-diff-finder/cache"
	cachemocks "file-diff-finder/cache/mocks"
	df "file-diff-finder/diff_finder"
	dfmocks "file-diff-finder/diff_finder/mocks"
	"file-diff-finder/router"
)

var _ = Describe("file diff endpoint API test - GET /diff", func() {
	Context("error cases", func() {
		var (
			mockCtrl      *gomock.Controller
			logger        *logrus.Logger
			cacheInstance *cachemocks.MockICache
			fdf           *dfmocks.MockIFileDiffFinder
		)
		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			logger = logrus.New()
			cacheInstance = cachemocks.NewMockICache(mockCtrl)
			fdf = dfmocks.NewMockIFileDiffFinder(mockCtrl)
		})

		AfterEach(func() {
			defer mockCtrl.Finish()
		})

		It("should return error due to request body parsing error", func() {
			req := &http.Request{
				URL:    &url.URL{},
				Method: http.MethodGet,
				Body:   http.NoBody,
			}

			rtr := router.NewRouter(logger, cacheInstance, fdf)
			res := http.HandlerFunc(rtr.Diff)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, req)

			decoder := json.NewDecoder(rr.Body)
			var t router.ErrorResp
			_ = decoder.Decode(&t)
			Expect(t.Code).To(BeEquivalentTo(http.StatusBadRequest))
			Expect(t.Message).To(BeEquivalentTo("error on parsing payload"))
		})

		It("should not validate negative version number", func() {
			text := "testetst"
			rp := &router.RequestPayload{
				Text:    &text,
				Version: -1,
			}

			payload, _ := json.Marshal(rp)

			req := &http.Request{
				URL:    &url.URL{},
				Method: http.MethodGet,
				Body:   io.NopCloser(bytes.NewReader(payload)),
			}

			rtr := router.NewRouter(logger, cacheInstance, fdf)
			res := http.HandlerFunc(rtr.Diff)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, req)

			decoder := json.NewDecoder(rr.Body)
			var t router.ErrorResp
			_ = decoder.Decode(&t)
			Expect(t.Code).To(BeEquivalentTo(http.StatusBadRequest))
			Expect(t.Message).To(BeEquivalentTo("Version must be provided and should be positive integer"))
		})

		It("should not validate negative version number", func() {
			rp := &router.RequestPayload{
				Version: 14,
			}

			payload, _ := json.Marshal(rp)

			req := &http.Request{
				URL:    &url.URL{},
				Method: http.MethodGet,
				Body:   io.NopCloser(bytes.NewReader(payload)),
			}

			rtr := router.NewRouter(logger, cacheInstance, fdf)
			res := http.HandlerFunc(rtr.Diff)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, req)

			decoder := json.NewDecoder(rr.Body)
			var t router.ErrorResp
			_ = decoder.Decode(&t)
			Expect(t.Code).To(BeEquivalentTo(http.StatusBadRequest))
			Expect(t.Message).To(BeEquivalentTo("missing text in payload"))
		})

		It("should not validate if file diff finder does not validate", func() {
			text := "testetst"
			rp := &router.RequestPayload{
				Text:    &text,
				Version: 12,
			}

			payload, _ := json.Marshal(rp)

			req := &http.Request{
				URL:    &url.URL{},
				Method: http.MethodGet,
				Body:   io.NopCloser(bytes.NewReader(payload)),
			}

			fdf.EXPECT().ValidateVersion(12).Return(errors.New("some error"))
			rtr := router.NewRouter(logger, cacheInstance, fdf)
			res := http.HandlerFunc(rtr.Diff)
			rr := httptest.NewRecorder()
			res.ServeHTTP(rr, req)

			decoder := json.NewDecoder(rr.Body)
			var t router.ErrorResp
			_ = decoder.Decode(&t)
			Expect(t.Code).To(BeEquivalentTo(http.StatusBadRequest))
			Expect(t.Message).To(BeEquivalentTo("some error"))
		})

		It("should reject request due to another running process", func() {
			text := "testetst"
			rp := &router.RequestPayload{
				Text:    &text,
				Version: 14,
			}

			payload, _ := json.Marshal(rp)

			cacheInstance.EXPECT().Get("14").Return(nil, true)
			req := &http.Request{
				URL:    &url.URL{},
				Method: http.MethodGet,
				Body:   io.NopCloser(bytes.NewReader(payload)),
			}

			fdf.EXPECT().ValidateVersion(14).Return(nil)
			rtr := router.NewRouter(logger, cacheInstance, fdf)
			handler := http.HandlerFunc(rtr.Diff)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			decoder := json.NewDecoder(rr.Body)
			var t router.ErrorResp
			_ = decoder.Decode(&t)
			Expect(t.Code).To(BeEquivalentTo(http.StatusUnprocessableEntity))
			Expect(t.Message).To(BeEquivalentTo("another process is getting the diff"))
		})
	})

	Context("response from diff function", func() {
		var (
			mockCtrl      *gomock.Controller
			logger        *logrus.Logger
			cacheInstance *cachemocks.MockICache
			fdf           *dfmocks.MockIFileDiffFinder
		)

		BeforeEach(func() {
			mockCtrl = gomock.NewController(GinkgoT())
			logger = logrus.New()
			cacheInstance = cachemocks.NewMockICache(mockCtrl)
			fdf = dfmocks.NewMockIFileDiffFinder(mockCtrl)
		})

		AfterEach(func() {
			defer mockCtrl.Finish()
		})

		It("should return delta", func() {
			text := "1bcabcabca6cabcab4abc"
			rp := &router.RequestPayload{
				Text:    &text,
				Version: 14,
			}

			payload, _ := json.Marshal(rp)

			cacheInstance.EXPECT().Get("14").Return(nil, false)
			cacheInstance.EXPECT().Set("14", true)
			cacheInstance.EXPECT().Delete("14")

			updatedIndeces := []df.UpdatedIndex{
				{
					OldValue: "b",
					NewValue: "a",
					Index:    10,
					Type:     df.ChangeTypesUpdated,
				},
			}
			fdf.EXPECT().ValidateVersion(14).Return(nil)
			fdf.EXPECT().Version().Return(13)
			fdf.EXPECT().Diff(gomock.Any(), "1bcabcabca6cabcab4abc").Return(updatedIndeces, nil)
			req := &http.Request{
				URL:    &url.URL{},
				Method: http.MethodGet,
				Body:   io.NopCloser(bytes.NewReader(payload)),
			}

			rtr := router.NewRouter(logger, cacheInstance, fdf)
			handler := http.HandlerFunc(rtr.Diff)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			decoder := json.NewDecoder(rr.Body)
			var t router.Response
			_ = decoder.Decode(&t)
			Expect(t.CurrentVersion).To(BeEquivalentTo(13))
			Expect(t.UpdatedVersion).To(BeEquivalentTo(14))
			Expect(t.Delta).To(BeEquivalentTo(updatedIndeces))
		})

		It("should return error", func() {
			text := "1bcabcabca6cabcab4abc"
			rp := &router.RequestPayload{
				Text:    &text,
				Version: 14,
			}

			payload, _ := json.Marshal(rp)

			cacheInstance.EXPECT().Get("14").Return(nil, false)
			cacheInstance.EXPECT().Set("14", true)
			cacheInstance.EXPECT().Delete("14")

			fdf.EXPECT().ValidateVersion(14).Return(nil)
			fdf.EXPECT().Diff(gomock.Any(), "1bcabcabca6cabcab4abc").Return(nil, errors.New("some error"))
			req := &http.Request{
				URL:    &url.URL{},
				Method: http.MethodGet,
				Body:   io.NopCloser(bytes.NewReader(payload)),
			}

			rtr := router.NewRouter(logger, cacheInstance, fdf)
			handler := http.HandlerFunc(rtr.Diff)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			decoder := json.NewDecoder(rr.Body)
			var t router.ErrorResp
			_ = decoder.Decode(&t)
			Expect(t.Code).To(BeEquivalentTo(http.StatusInternalServerError))
			Expect(t.Message).To(BeEquivalentTo("some error"))
		})
	})
})

var _ = Describe("liveness endpoint API test - GET /live", func() {
	var (
		mockCtrl      *gomock.Controller
		logger        *logrus.Logger
		cacheInstance *cachemocks.MockICache
		fdf           *dfmocks.MockIFileDiffFinder
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger = logrus.New()
		cacheInstance = cachemocks.NewMockICache(mockCtrl)
		fdf = dfmocks.NewMockIFileDiffFinder(mockCtrl)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	It("should return ok", func() {
		req := &http.Request{
			URL:    &url.URL{},
			Method: http.MethodGet,
		}

		rtr := router.NewRouter(logger, cacheInstance, fdf)
		res := http.HandlerFunc(rtr.Liveness)
		rr := httptest.NewRecorder()
		res.ServeHTTP(rr, req)

		Expect(rr.Code).To(BeEquivalentTo(http.StatusOK))
	})
})

var _ = Describe("readiness endpoint API test - GET /readiness", func() {
	var (
		mockCtrl      *gomock.Controller
		logger        *logrus.Logger
		cacheInstance *cachemocks.MockICache
		fdf           *dfmocks.MockIFileDiffFinder
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		logger = logrus.New()
		cacheInstance = cachemocks.NewMockICache(mockCtrl)
		fdf = dfmocks.NewMockIFileDiffFinder(mockCtrl)
	})

	AfterEach(func() {
		defer mockCtrl.Finish()
	})

	It("should return ok if cache returns response", func() {
		req := &http.Request{
			URL:    &url.URL{},
			Method: http.MethodGet,
		}
		cacheInstance.EXPECT().Get(cache.ReadinessKey).Return(nil, true)

		rtr := router.NewRouter(logger, cacheInstance, fdf)
		res := http.HandlerFunc(rtr.Readiness)
		rr := httptest.NewRecorder()
		res.ServeHTTP(rr, req)

		Expect(rr.Code).To(BeEquivalentTo(http.StatusOK))
	})

	It("should not return ok unless cache returns response", func() {
		req := &http.Request{
			URL:    &url.URL{},
			Method: http.MethodGet,
		}
		cacheInstance.EXPECT().Get(cache.ReadinessKey).Return(nil, false)

		rtr := router.NewRouter(logger, cacheInstance, fdf)
		res := http.HandlerFunc(rtr.Readiness)
		rr := httptest.NewRecorder()
		res.ServeHTTP(rr, req)

		decoder := json.NewDecoder(rr.Body)
		var t router.ErrorResp
		_ = decoder.Decode(&t)
		Expect(t.Code).To(BeEquivalentTo(http.StatusInternalServerError))
		Expect(t.Message).To(BeEquivalentTo("error cache is not ok"))
	})
})
