package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hugosrc/shortlink/internal/core/domain"
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
	originalURL := "http://github.com/hugosrc/shortlink"

	mockRepository.EXPECT().FindByHash(context.Background(), hash).Return(&domain.Link{
		Hash:         hash,
		OriginalURL:  originalURL,
		UserID:       "88b144af-7743-4824-a3b4-01839600bcbb",
		CreationTime: time.Now(),
	}, nil)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	url, err := svc.FindByHash(context.Background(), hash)

	assert.Equal(t, originalURL, url)
	assert.Nil(t, err)
}

func TestLinkService_FindUrlByHash_RepositoryErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(nil, assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	url, err := svc.FindByHash(context.Background(), "a74dAS1")

	assert.Empty(t, url)
	assert.NotNil(t, err)
}

func TestLinkService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	originalURL := "http://github.com/hugosrc/shortlink"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(&domain.Link{
		Hash:         hash,
		OriginalURL:  originalURL,
		UserID:       userId,
		CreationTime: time.Now(),
	}, nil)

	mockRepository.EXPECT().Delete(context.Background(), gomock.Any()).Return(nil)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	err := svc.Delete(context.Background(), hash, userId)

	assert.Nil(t, err)
}

func TestLinkService_Delete_RepositoryFindErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(nil, assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	err := svc.Delete(context.Background(), hash, userId)

	assert.NotNil(t, err)
}

func TestLinkService_Delete_UserErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	originalURL := "http://github.com/hugosrc/shortlink"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"
	invalidUserID := "abae413b-1f02-4e8c-8895-3ab56f46b651"

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(&domain.Link{
		Hash:         hash,
		OriginalURL:  originalURL,
		UserID:       userId,
		CreationTime: time.Now(),
	}, nil)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	err := svc.Delete(context.Background(), hash, invalidUserID)

	assert.NotNil(t, err)
}

func TestLinkService_Delete_RepositoryDeleteErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	originalURL := "http://github.com/hugosrc/shortlink"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(&domain.Link{
		Hash:         hash,
		OriginalURL:  originalURL,
		UserID:       userId,
		CreationTime: time.Now(),
	}, nil)
	mockRepository.EXPECT().Delete(context.Background(), gomock.Any()).Return(assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	err := svc.Delete(context.Background(), hash, userId)

	assert.NotNil(t, err)
}

func TestLinkService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	originalURL := "http://github.com/hugosrc/shortlink"
	newURL := "http://github.com/hugosrc/surf-forecast-api"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"
	creationTime := time.Now()

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(&domain.Link{
		Hash:         hash,
		OriginalURL:  originalURL,
		UserID:       userId,
		CreationTime: creationTime,
	}, nil)

	mockRepository.EXPECT().Update(context.Background(), gomock.Any(), gomock.Any()).Return(nil)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	link, err := svc.Update(context.Background(), hash, newURL, userId)

	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, link.Hash, hash)
	assert.Equal(t, link.OriginalURL, newURL)
	assert.Equal(t, link.UserID, userId)
	assert.Equal(t, link.CreationTime, creationTime)
}

func TestLinkService_Update_RepositoryFindErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	newURL := "http://github.com/hugosrc/surf-forecast-api"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(nil, assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	link, err := svc.Update(context.Background(), hash, newURL, userId)

	assert.NotNil(t, err)
	assert.Nil(t, link)
}

func TestLinkService_Update_UserErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	originalURL := "http://github.com/hugosrc/shortlink"
	newURL := "http://github.com/hugosrc/surf-forecast-api"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"
	invalidUserID := "abae413b-1f02-4e8c-8895-3ab56f46b651"

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(&domain.Link{
		Hash:         hash,
		OriginalURL:  originalURL,
		UserID:       userId,
		CreationTime: time.Now(),
	}, nil)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	link, err := svc.Update(context.Background(), hash, newURL, invalidUserID)

	assert.NotNil(t, err)
	assert.Nil(t, link)
}

func TestLinkService_Update_RepositoryUpdateErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCounter := mock.NewMockCounter(ctrl)
	mockEncoder := mock.NewMockEncoder(ctrl)
	mockRepository := mock.NewMockLinkRepository(ctrl)

	hash := "15saSAE"
	originalURL := "http://github.com/hugosrc/shortlink"
	newURL := "http://github.com/hugosrc/surf-forecast-api"
	userId := "88b144af-7743-4824-a3b4-01839600bcbb"

	mockRepository.EXPECT().FindByHash(context.Background(), gomock.Any()).Return(&domain.Link{
		Hash:         hash,
		OriginalURL:  originalURL,
		UserID:       userId,
		CreationTime: time.Now(),
	}, nil)
	mockRepository.EXPECT().Update(context.Background(), gomock.Any(), gomock.Any()).Return(assert.AnError)

	svc := service.NewLinkService(mockCounter, mockEncoder, mockRepository)
	link, err := svc.Update(context.Background(), hash, newURL, userId)

	assert.NotNil(t, err)
	assert.Nil(t, link)
}
