package swim

import "github.com/ApogeeNetworking/dnacenter/requests"

// DNA-C SoftWare Image Management (SWIM)

// Service ...
type Service struct {
	baseURL string
	http    *requests.Req
}

// New creates an instance of a DNA-C SWIM Service
func New(uri string, r *requests.Req) *Service {
	return &Service{baseURL: uri + "/intent/api/v1/image", http: r}
}

// Image ... DNA-C NetDevice Image
type Image struct {
	ImageUUID                 string   `json:"imageUuid"`
	Name                      string   `json:"name"`
	Family                    string   `json:"family"`
	Version                   string   `json:"version"`
	Md5CheckSum               string   `json:"md5Checksum"`
	ShaCheckSum               string   `json:"shaCheckSum"`
	CreatedTime               string   `json:"createdTime"`
	ImageType                 string   `json:"imageType"`
	FileSize                  string   `json:"fileSize"`
	ImageName                 string   `json:"imageName"`
	ApplicationType           string   `json:"applicationType"`
	Feature                   string   `json:"feature"`
	FileServiceID             string   `json:"fileServiceId"`
	IsTaggedGolden            bool     `json:"isTaggedGolden"`
	ImageSeries               []string `json:"imageSeries"`
	ImageSource               string   `json:"imageSource"`
	ExtendedAttributes        struct{} `json:"extendedAttributes"`
	Vendor                    string   `json:"vendor"`
	ImageIntegrityStatus      string   `json:"imageIntegrityStatus"`
	ApplicableDevicesForImage []struct {
		MdfID       string   `json:"mdfId"`
		ProductName string   `json:"productName"`
		ProductID   []string `json:"productId"`
	} `json:"applicableDevicesForImage"`
	ImportSourceType string `json:"importSourceType"`
	CCOReverseSync   bool   `json:"ccoreverseSync"`
}
