package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hugosrc/shortlink/internal/core/service"
	"github.com/hugosrc/shortlink/internal/core/service/mock"
	"github.com/stretchr/testify/assert"
)

func TestLinkService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "5s48ASE"
	url := "http://github.com/hugosrc/shortlink"
	userID := "88b144af-7743-4824-a3b4-01839600bcbb"

	mockCounter.EXPECT().Inc().Return(10000, nil)
	mockEncoder.EXPECT().EncodeToString(gomock.Any()).Return(hash)
	mockRepository.EXPECT().Create(context.Background(), gomock.Any()).Return(nil)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)

	link, err := svc.Create(context.Background(), url, userID)

	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, link.Hash, hash)
	assert.Equal(t, link.OriginalURL, url)
	assert.Equal(t, link.UserID, userID)
}

func TestLinkService_Create_CounterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	mockCounter.EXPECT().Inc().Return(0, assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)

	link, err := svc.Create(
		context.Background(),
		"http://github.com/hugosrc/shortlink",
		"88b144af-7743-4824-a3b4-01839600bcbb")

	assert.Nil(t, link)
	assert.NotNil(t, err)
}

func TestLinkService_Create_RepositoryErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	mockCounter.EXPECT().Inc().Return(10000, nil)
	mockEncoder.EXPECT().EncodeToString(gomock.Any()).Return("as544ca")
	mockRepository.EXPECT().Create(context.Background(), gomock.Any()).Return(assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)

	link, err := svc.Create(
		context.Background(),
		"http://github.com/hugosrc/shortlink",
		"88b144af-7743-4824-a3b4-01839600bcbb")

	assert.Nil(t, link)
	assert.NotNil(t, err)
}

func TestLinkService_FindUrlByHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "r74aASW"
	originalURL := "88b144af-7743-4824-a3b4-01839600bcbb"

	mockRepository.EXPECT().FindUrlByHash(context.Background(), hash).Return(originalURL, nil)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	url, err := svc.FindUrlByHash(context.Background(), hash)

	assert.Equal(t, originalURL, url)
	assert.Nil(t, err)
}

func TestLinkService_FindUrlByHash_RepositoryErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	mockRepository.EXPECT().FindUrlByHash(context.Background(), gomock.Any()).Return("", assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	url, err := svc.FindUrlByHash(context.Background(), "a74dAS1")

	assert.Empty(t, url)
	assert.NotNil(t, err)
}
