// mempool/mempool.go
package mempool

import (
	"doc-tracker/models"
	"doc-tracker/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gofiber/fiber"
)

type TrackerEntry = models.Tracker

var (
	mempool      = make(map[string]*models.Tracker)
	mu           sync.RWMutex
	jsonFilePath = "data/mempool.json"
	binFilePath  = "data/mempool.bin"
	pubKeyPath   = "data/public.key"
	privKeyPath  = "data/private.key"
)

// InitKeys inisialisasi atau load kunci
func InitKeys() (*utils.ECDHKeyPair, error) {
	// Jika kunci sudah ada, load dari file
	if _, err := os.Stat(pubKeyPath); err == nil {
		return utils.LoadKeys()
	}

	// Generate new key pair
	keyPair, err := utils.GenerateECDHKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate keys: %v", err)
	}

	// Simpan public key
	pubKeyBytes := utils.SerializePublicKey(keyPair.PublicKey)
	if err := os.WriteFile(pubKeyPath, pubKeyBytes, 0600); err != nil {
		return nil, fmt.Errorf("failed to save public key: %v", err)
	}

	// Simpan private key
	privKeyBytes := keyPair.PrivateKey.Bytes()
	if err := os.WriteFile(privKeyPath, privKeyBytes, 0600); err != nil {
		return nil, fmt.Errorf("failed to save private key: %v", err)
	}

	return keyPair, nil
}

// InitEncryptMempool migrasi ke mempool terenkripsi
func InitEncryptMempool() error {
	mu.Lock()
	defer mu.Unlock()

	// Skip jika file terenkripsi sudah ada
	if _, err := os.Stat(binFilePath); err == nil {
		log.Println("🟡 Encrypted mempool already exists")
		return nil
	}

	// Cek jika file plaintext ada
	if _, err := os.Stat(jsonFilePath); os.IsNotExist(err) {
		log.Println("🟡 No mempool found to encrypt")
		return nil
	}

	log.Println("🔄 Encrypting mempool...")

	// Load kunci
	keyPair, err := utils.LoadKeys()
	if err != nil {
		return fmt.Errorf("failed to load keys: %v", err)
	}

	// Baca data plaintext
	plaintext, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read mempool: %v", err)
	}

	// Enkripsi data
	// Untuk demo, kita gunakan kunci public sendiri sebagai peer
	sharedSecret, err := utils.DeriveSharedSecret(keyPair.PrivateKey, keyPair.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to derive secret: %v", err)
	}

	encKey := utils.GenerateEncryptionKey(sharedSecret)
	ciphertext, err := utils.EncryptData(encKey, plaintext)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// Simpan versi terenkripsi
	if err := os.WriteFile(binFilePath, ciphertext, 0600); err != nil {
		return fmt.Errorf("failed to save encrypted mempool: %v", err)
	}

	// Hapus file plaintext (opsional)
	if err := os.Remove(jsonFilePath); err != nil {
		log.Printf("⚠️ Could not remove plaintext file: %v", err)
	}

	log.Println("✅ Mempool encrypted successfully")
	return nil
}

// LoadFromFile memuat mempool dari file terenkripsi
func LoadFromFile() error {
	mu.Lock()
	defer mu.Unlock()

	// Coba load dari file terenkripsi dulu
	if _, err := os.Stat(binFilePath); err == nil {
		// Load kunci
		keyPair, err := utils.LoadKeys()
		if err != nil {
			return fmt.Errorf("failed to load keys: %v", err)
		}

		// Baca data terenkripsi
		ciphertext, err := os.ReadFile(binFilePath)
		if err != nil {
			return fmt.Errorf("failed to read encrypted mempool: %v", err)
		}

		// Dekripsi data
		sharedSecret, err := utils.DeriveSharedSecret(keyPair.PrivateKey, keyPair.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to derive secret: %v", err)
		}

		encKey := utils.GenerateEncryptionKey(sharedSecret)
		plaintext, err := utils.DecryptData(encKey, ciphertext)
		if err != nil {
			return fmt.Errorf("decryption failed: %v", err)
		}

		// Parse JSON
		var loadedMempool map[string]*models.Tracker
		if err := json.Unmarshal(plaintext, &loadedMempool); err != nil {
			return fmt.Errorf("failed to parse mempool: %v", err)
		}

		mempool = loadedMempool
		log.Println("✅ Mempool loaded from encrypted file")
		return nil
	}

	// Fallback ke file plaintext
	if _, err := os.Stat(jsonFilePath); err == nil {
		log.Println("🟡 Loading from plaintext mempool")
		return loadFromPlaintext()
	}

	return errors.New("no mempool file found")
}

// loadFromPlaintext helper untuk load dari JSON plaintext
func loadFromPlaintext() error {
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read plaintext mempool: %v", err)
	}

	var loadedMempool map[string]*models.Tracker
	if err := json.Unmarshal(data, &loadedMempool); err != nil {
		return fmt.Errorf("failed to parse mempool: %v", err)
	}

	mempool = loadedMempool
	return nil
}

// SaveToFile menyimpan mempool ke file terenkripsi
func SaveToFile() error {
	mu.Lock()
	defer mu.Unlock()

	// Marshal data ke JSON
	data, err := json.MarshalIndent(mempool, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal mempool: %v", err)
	}

	// Load kunci
	keyPair, err := utils.LoadKeys()
	if err != nil {
		return fmt.Errorf("failed to load keys: %v", err)
	}

	// Enkripsi data
	sharedSecret, err := utils.DeriveSharedSecret(keyPair.PrivateKey, keyPair.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to derive secret: %v", err)
	}

	encKey := utils.GenerateEncryptionKey(sharedSecret)
	ciphertext, err := utils.EncryptData(encKey, data)
	if err != nil {
		return fmt.Errorf("encryption failed: %v", err)
	}

	// Simpan ke file
	if err := os.WriteFile(binFilePath, ciphertext, 0600); err != nil {
		return fmt.Errorf("failed to save encrypted mempool: %v", err)
	}

	return nil
}

// [Fungsi-fungsi manajemen mempool yang sama seperti sebelumnya...]
// Add, GetAll, GetCompletedTrackers, RemoveFromMempool, dll.

// Add tracker ke mempool jika belum ada
func Add(t *models.Tracker) {
	if _, exists := mempool[t.ID]; !exists {
		mempool[t.ID] = t
	}
	SaveToFile()
}

// Get semua tracker di mempool
func GetAll() []*models.Tracker {
	var list []*models.Tracker
	for _, t := range mempool {
		list = append(list, t)
	}
	return list
}

func GetByID(id string) *models.Tracker {
	mu.RLock()
	defer mu.RUnlock()

	if tracker, exists := mempool[id]; exists {
		return tracker
	}
	return nil
}

// Ambil tracker dengan status complete
func GetCompletedTrackers() []*models.Tracker {
	var completed []*models.Tracker
	for _, t := range mempool {
		if t.Status == "complete" {
			completed = append(completed, t)
		}
	}
	return completed
}

// Hapus tracker dari mempool
func RemoveFromMempool(id string) {
	delete(mempool, id)
}

func SyncMempool(c *fiber.Ctx) error {
	var tracker models.Tracker
	if err := c.BodyParser(&tracker); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid tracker data"})
	}

	Add(&tracker) // gunakan fungsi baru
	return c.JSON(fiber.Map{"status": "tracker synced"})
}

func GetProgressTrackers() []*models.Tracker {
	var list []*models.Tracker
	for _, t := range mempool {
		if t.Status == "progress" {
			list = append(list, t)
		}
	}
	return list
}

func Update(tracker models.Tracker) {
	if _, exists := mempool[tracker.ID]; exists {
		mempool[tracker.ID] = &tracker
		return
	}
	// jika tidak ditemukan, tambahkan
	mempool[tracker.ID] = &tracker
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()
	mempool = make(map[string]*models.Tracker)
}

func AddIfNotExists(tracker models.Tracker) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := mempool[tracker.ID]; !exists {
		mempool[tracker.ID] = &tracker
	}
}

func Exists(id string) bool {
	mu.RLock()
	defer mu.RUnlock()

	_, exists := mempool[id]
	return exists
}

// Iterate iterates over all transactions in the mempool and applies the given function.
func Iterate(fn func(tx *models.Tracker) error) error {

	for _, tx := range mempool { // assuming mempool is a slice or map of *models.Transaction
		if err := fn(tx); err != nil {
			return err
		}
	}
	return nil
}

func UpdateTracker(tracker *models.Tracker) {
	mu.Lock()
	defer mu.Unlock()

	mempool[tracker.ID] = tracker

	if err := SaveToFile(); err != nil {
		fmt.Printf("❌ Gagal simpan mempool: %v\n", err)
	}
}

func RemoveDuplicateEntries() {
	unique := make(map[string]bool)
	var cleaned []TrackerEntry

	for _, entry := range mempool {
		if !unique[entry.ID] {
			unique[entry.ID] = true
			cleaned = append(cleaned, *entry)
		}
	}

	// Rebuild mempool with unique entries
	mempool = make(map[string]*models.Tracker)
	for i := range cleaned {
		mempool[cleaned[i].ID] = &cleaned[i]
	}
	fmt.Printf("[Cleanse] %d unique tracker entries retained\n", len(cleaned))
}
