package services

import (
	"doc-tracker/blockchain"
	"doc-tracker/mempool"
	"doc-tracker/models"
	"doc-tracker/p2p"
	"doc-tracker/utils"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	eciesgo "github.com/ecies/go/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var trackerMap = map[string]models.Tracker{}

const trackerFile = "data/trackers.json"

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	secret := os.Getenv("TRACKER_SECRET")
	if secret == "" {
		panic("TRACKER_SECRET tidak di-set di environment. Pastikan .env dibaca.")
	}
	secretBytes, _ := hex.DecodeString(secret)
	if len(secretBytes) != 32 {
		panic(fmt.Sprintf("TRACKER_SECRET harus 32 byte (256 bit) untuk ECIES Private Key, tetapi mendapatkan %d byte", len(secretBytes)))
	}
	priv := eciesgo.NewPrivateKeyFromBytes(secretBytes) // generate new key
	pub := priv.PublicKey
	fmt.Println("Private:", priv.Hex())
	fmt.Println("Public:", pub.Hex(true))

	// Konversi ke *ecdsa.PrivateKey dan *ecdsa.PublicKey
	// Konversi ECIES ke ECDSA
	privECDSA := utils.ECIESPrivateKeyToECDSA(priv)
	pubECDSA := &privECDSA.PublicKey

	// Simpan ke file .pem
	if err := utils.SavePEMKey("data/private.pem", privECDSA); err != nil {
		panic(err)
	}
	if err := utils.SavePEMPub("data/public.pem", pubECDSA); err != nil {
		panic(err)
	}

	fmt.Println("Kunci berhasil disimpan ke private.pem dan public.pem")

	LoadTrackersFromDisk()
}

func CreateTracker(input models.Tracker) (models.Tracker, error) {
	input.ID = uuid.New().String()
	input.CreatedAt = time.Now().Unix()
	input.Status = "progress"

	// Generate wallet/address untuk pengaju
	senderWallet := GetOrCreateWallet(input.Creator)
	input.CreatorAddr = senderWallet.Address

	// Enkripsi checkpoint jika perlu
	for i, cp := range input.Checkpoints {
		receiverWallet := GetOrCreateWallet(cp.Email)
		input.Checkpoints[i].Address = receiverWallet.Address

		// Jika boleh melihat isi dokumen, enkripsi
		if cp.IsViewable {
			encrypted := utils.EncryptWithPublicKey(cp.Note, receiverWallet.PublicKey)
			input.Checkpoints[i].EncryptedNote = encrypted
			input.Checkpoints[i].Note = "" // kosongkan untuk keamanan
		} else {
			// Tidak boleh melihat, kosongkan isinya
			input.Checkpoints[i].Note = ""
			input.Checkpoints[i].EncryptedNote = ""
		}
	}

	// Simpan ke mempool dan broadcast
	mempool.Add(&input)
	p2p.BroadcastToMempool(input)

	return input, nil
}

func GetDataTracker() ([]models.Tracker, error) {
	var trackers []models.Tracker
	// Get from mempool
	err := mempool.Iterate(func(tx *models.Tracker) error {
		trackers = append(trackers, *tx)
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Get from blockchain
	errBlc := blockchain.Iterate(func(tx *models.Tracker) error {
		trackers = append(trackers, *tx)
		return nil
	})
	if errBlc != nil {
		return nil, errBlc
	}

	return trackers, nil
}

func GetTrackers() ([]models.Tracker, error) {
	trackers, err := GetDataTracker()
	if err != nil {
		return nil, err
	}
	return trackers, nil
}
func GetTrackerByID(id string) (models.Tracker, error) {
	trackers, err := GetDataTracker()
	if err != nil {
		return models.Tracker{}, err
	}
	for _, t := range trackers {
		if t.ID == id {
			return t, nil
		}
	}
	return models.Tracker{}, utils.ErrNotFound
}

func GetTrackersByAddress(address string) ([]models.Tracker, error) {
	trackers, err := GetDataTracker()
	if err != nil {
		return nil, err
	}
	for i := len(trackers) - 1; i >= 0; i-- {
		if !hasCheckpointForAddress(&trackers[i], address) {
			// Hapus tracker jika tidak ada checkpoint untuk address ini
			trackers = append(trackers[:i], trackers[i+1:]...)
		}
	}
	if len(trackers) == 0 {
		return nil, utils.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return trackers, nil
}

func GetTrackerByHash(hash string) (models.Tracker, error) {
	trackers, err := GetDataTracker()
	if err != nil {
		return models.Tracker{}, err
	}

	for _, t := range trackers {
		for _, cp := range t.Checkpoints {
			if cp.EvidenceHash == hash {
				return t, nil
			}
		}
	}
	return models.Tracker{}, utils.ErrNotFound
}

func GetEvidencePath(tracker models.Tracker, hash string) string {
	// Implement logic to get the evidence path based on trackerID and hash
	for _, cp := range tracker.Checkpoints {
		if cp.EvidencePath != "" && cp.EvidenceHash == hash {
			return cp.EvidencePath
		}
	}
	return ""
}

func hasCheckpointForAddress(tx *models.Tracker, address string) bool {
	for _, cp := range tx.Checkpoints {
		if cp.Address == address {
			return true
		}
	}
	return false
}

func SaveTrackersToDisk() {
	utils.SaveToFile(trackerFile, trackerMap)
}

func LoadTrackersFromDisk() {
	utils.LoadFromFile(trackerFile, &trackerMap)
}

func GetTrackerSummary(email string) (map[string]int, error) {
	trackers, err := GetDataTracker()
	if err != nil {
		return nil, err
	}

	summary := make(map[string]int)
	for _, tracker := range trackers {
		if tracker.Creator == email {
			summary[tracker.Status]++ // Tambahkan status tracker yang dibuat oleh email ini
			continue                  // Lanjutkan ke tracker berikutnya jika ini adalah tracker yang dibuat oleh email
		}
		for _, cp := range tracker.Checkpoints {
			if cp.Email == email {
				summary[tracker.Status]++
			}
		}
	}

	return summary, nil
}
