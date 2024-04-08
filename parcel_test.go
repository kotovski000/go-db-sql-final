package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	// настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Client, parcel.Client)
	assert.Equal(t, res.Status, parcel.Status)
	assert.Equal(t, res.Address, parcel.Address)
	assert.Equal(t, res.CreatedAt, parcel.CreatedAt)
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
	err = store.Delete(id)
	require.NoError(t, err)

	res, err = store.Get(parcel.Number)
	require.Error(t, err, sql.ErrNoRows)
	require.Empty(t, res)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	// настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)
	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что адрес обновился
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Address, "new test address")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	// настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()
	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)
	// set status
	// обновите статус, убедитесь в отсутствии ошибки
	err = store.SetStatus(id, ParcelStatusDelivered)
	require.NoError(t, err)
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
	res, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, res.Status, ParcelStatusDelivered)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	// настройте подключение к БД
	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i]) // добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
		require.NoError(t, err)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client) // получите список посылок по идентификатору клиента, сохранённого в переменной client
	// убедитесь в отсутствии ошибки
	require.NoError(t, err)
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	require.Equal(t, len(parcels), len(storedParcels))
	// check
	for _, parcel := range storedParcels {
		// в parcelMap лежат добавленные посылки, ключ - идентификатор посылки, значение - сама посылка
		originalParcel, ok := parcelMap[parcel.Number]
		// убедитесь, что все посылки из storedParcels есть в parcelMap
		require.True(t, ok)
		// убедитесь, что значения полей полученных посылок заполнены верно
		require.Equal(t, originalParcel, parcel)
	}
}
