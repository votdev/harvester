package downloadurl

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"

	harvesterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	"github.com/harvester/harvester/pkg/config"
	ctlharvesterv1 "github.com/harvester/harvester/pkg/generated/controllers/harvesterhci.io/v1beta1"
)

const (
	controllerName = "downloadurl-controller"
)

func Register(ctx context.Context, management *config.Management, options config.Options) error {
	downloadURLController := management.HarvesterFactory.Harvesterhci().V1beta1().DownloadURL()
	handler := &handler{
		ctl: downloadURLController,
		httpClient: http.Client{
			Timeout: 15 * time.Second,
		},
	}
	downloadURLController.OnChange(ctx, controllerName, handler.OnChange)
	return nil
}

// handler tries to download from specified URL on DownloadURL changes.
type handler struct {
	ctl        ctlharvesterv1.DownloadURLController
	httpClient http.Client
}

func (h *handler) OnChange(_ string, downloadURL *harvesterv1.DownloadURL) (*harvesterv1.DownloadURL, error) {
	if downloadURL == nil || downloadURL.DeletionTimestamp != nil {
		return nil, nil
	}

	downloadURLCopy := downloadURL.DeepCopy()

	switch downloadURL.Status.Status {
	case "":
		downloadURLCopy.Status.Status = harvesterv1.INIT_STATE
	case harvesterv1.INIT_STATE:
		downloadURLCopy.Status.Status = harvesterv1.DOWNLOADING_STATE
		logrus.Debugf("download from URL: %s", downloadURL.URL)
		resp, err := h.httpClient.Get(downloadURL.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		// then operates the response body
	case harvesterv1.DOWNLOADING_STATE:
		downloadURLCopy.Status.Status = harvesterv1.DONE_STATE
	}

	if !reflect.DeepEqual(downloadURLCopy, downloadURL) {
		return h.ctl.UpdateStatus(downloadURLCopy)
	}

	return downloadURL, nil
}
