// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tasks

import (
	"context"
	"fmt"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type BaremetalSyncConfigTask struct {
	SBaremetalBaseTask
}

func init() {
	taskman.RegisterTask(BaremetalSyncConfigTask{})
}

func (self *BaremetalSyncConfigTask) OnInit(ctx context.Context, obj db.IStandaloneModel, body jsonutils.JSONObject) {
	baremetal := obj.(*models.SHost)
	if baremetal.IsBaremetal {
		self.DoSyncConfig(ctx, baremetal)
	} else {
		self.SetStageComplete(ctx, nil)
	}
}

func (self *BaremetalSyncConfigTask) DoSyncConfig(ctx context.Context, baremetal *models.SHost) {
	self.SetStage("OnSyncConfigComplete", nil)
	url := fmt.Sprintf("/baremetals/%s/sync-config", baremetal.Id)
	headers := self.GetTaskRequestHeader()
	_, err := baremetal.BaremetalSyncRequest(ctx, "POST", url, headers, nil)
	if err != nil {
		self.SetStageFailed(ctx, err.Error())
	}
}

func (self *BaremetalSyncConfigTask) OnSyncConfigComplete(ctx context.Context, baremetal *models.SHost, body jsonutils.JSONObject) {
	self.SetStage("OnSyncstatusComplete", nil)
	baremetal.StartSyncstatus(ctx, self.UserCred, self.GetTaskId())
	logclient.AddActionLogWithStartable(self, baremetal, logclient.ACT_SYNC_CONF, "", self.UserCred, true)
}

func (self *BaremetalSyncConfigTask) OnSyncstatusComplete(ctx context.Context, baremetal *models.SHost, body jsonutils.JSONObject) {
	self.SetStageComplete(ctx, nil)
}
