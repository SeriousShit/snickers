package downloaders

import (
	"strconv"
	"time"

	"../db"

	"code.cloudfoundry.org/lager"

	"github.com/cavaliercoder/grab"
	"github.com/flavioribeiro/gonfig"

)

// HTTPDownload function downloads sources using
// http protocol.
func HTTPDownload(logger lager.Logger, config gonfig.Gonfig, dbInstance db.Storage, jobID string) error {
	log := logger.Session("http-download")
	log.Info("start", lager.Data{"job": jobID})
	defer log.Info("finished")

	job, err := dbInstance.RetrieveJob(jobID)
	if err != nil {
		return err
	}


	client := grab.NewClient()
	req, err := grab.NewRequest(job.LocalSource, job.Source)

	resp := client.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			job, err = dbInstance.RetrieveJob(jobID)
			if err != nil {
				return err
			}

			percentage := strconv.FormatInt(int64(resp.Progress()), 10)

			if job.Progress != percentage {
				job.Progress = percentage + "%"
				dbInstance.UpdateJob(job.ID, job)
			}

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}


	//respch, err := grab.GetAsync(job.LocalSource, job.Source)
	if err != nil {
		return nil
	}

	return nil
}
