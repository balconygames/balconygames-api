package collector

import "gitlab.com/balconygames/analytics/pkg/runtime"

/*

collector:
	- subscribe to MSQ and reading the messages
	- in batches push messages to click house or other repositories
*/

func New(r *runtime.Runtime) error {
	return nil
}
