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

package guestdrivers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/osprofile"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/quotas"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/cloudprovider"
	guestdriver_types "yunion.io/x/onecloud/pkg/compute/guestdrivers/types"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/util/billing"
)

type SBaseGuestScheduleDriver struct{}

func (d SBaseGuestScheduleDriver) DoScheduleSKUFilter() bool { return true }

func (d SBaseGuestScheduleDriver) DoScheduleCPUFilter() bool { return true }

func (d SBaseGuestScheduleDriver) DoScheduleMemoryFilter() bool { return true }

func (d SBaseGuestScheduleDriver) DoScheduleStorageFilter() bool { return true }

func (d SBaseGuestScheduleDriver) DoScheduleCloudproviderTagFilter() bool { return false }

type SBaseGuestDriver struct {
	SBaseGuestScheduleDriver
}

func (self *SBaseGuestDriver) StartGuestCreateTask(guest *models.SGuest, ctx context.Context, userCred mcclient.TokenCredential, data *jsonutils.JSONDict, pendingUsage quotas.IQuota, parentTaskId string) error {
	task, err := taskman.TaskManager.NewTask(ctx, "GuestCreateTask", guest, userCred, data, parentTaskId, "", pendingUsage)
	if err != nil {
		return err
	}
	task.ScheduleRun(nil)
	return nil
}

func (self *SBaseGuestDriver) OnGuestCreateTaskComplete(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	duration, _ := task.GetParams().GetString("duration")
	if len(duration) > 0 {
		bc, err := billing.ParseBillingCycle(duration)
		if err == nil && guest.ExpiredAt.IsZero() {
			guest.SaveRenewInfo(ctx, task.GetUserCred(), &bc, nil, "")
		}
		if jsonutils.QueryBoolean(task.GetParams(), "auto_prepaid_recycle", false) {
			err := guest.CanPerformPrepaidRecycle()
			if err == nil {
				task.SetStageComplete(ctx, nil)
				guest.DoPerformPrepaidRecycle(ctx, task.GetUserCred(), true)
				return nil
			}
		}
	}
	if jsonutils.QueryBoolean(task.GetParams(), "auto_start", false) {
		task.SetStage("on_auto_start_guest", nil)
		return guest.StartGueststartTask(ctx, task.GetUserCred(), nil, task.GetTaskId())
	} else {
		task.SetStage("on_sync_status_complete", nil)
		return guest.StartSyncstatus(ctx, task.GetUserCred(), task.GetTaskId())
	}
}

func (self *SBaseGuestDriver) StartDeleteGuestTask(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, params *jsonutils.JSONDict, parentTaskId string) error {
	task, err := taskman.TaskManager.NewTask(ctx, "GuestDeleteTask", guest, userCred, params, parentTaskId, "", nil)
	if err != nil {
		return err
	}
	task.ScheduleRun(nil)
	return nil
}

func (self *SBaseGuestDriver) ValidateImage(ctx context.Context, image *cloudprovider.SImage) error {
	return nil
}

func (self *SBaseGuestDriver) RequestDetachDisksFromGuestForDelete(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	task.ScheduleRun(nil)
	return nil
}

func (self *SBaseGuestDriver) OnDeleteGuestFinalCleanup(ctx context.Context, guest *models.SGuest, userCred mcclient.TokenCredential) error {
	return guest.DeleteAllDisksInDB(ctx, userCred)
}

func (self *SBaseGuestDriver) RequestDetachDisk(ctx context.Context, guest *models.SGuest, disk *models.SDisk, task taskman.ITask) error {
	task.ScheduleRun(nil)
	return nil
}

func (self *SBaseGuestDriver) RequestAttachDisk(ctx context.Context, guest *models.SGuest, disk *models.SDisk, task taskman.ITask) error {
	task.ScheduleRun(nil)
	return nil
}

func (self *SBaseGuestDriver) RequestOpenForward(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, req *guestdriver_types.OpenForwardRequest) (*guestdriver_types.OpenForwardResponse, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SBaseGuestDriver) RequestListForward(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, req *guestdriver_types.ListForwardRequest) (*guestdriver_types.ListForwardResponse, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SBaseGuestDriver) RequestCloseForward(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, req *guestdriver_types.CloseForwardRequest) (*guestdriver_types.CloseForwardResponse, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SBaseGuestDriver) RequestSaveImage(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, task taskman.ITask) error {
	return errors.Wrapf(cloudprovider.ErrNotImplemented, "RequestSaveImage")
}

func (self *SBaseGuestDriver) RequestGuestCreateAllDisks(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) GetDetachDiskStatus() ([]string, error) {
	return []string{}, fmt.Errorf("This Guest driver dose not implement GetDetachDiskStatus")
}

func (self *SBaseGuestDriver) GetAttachDiskStatus() ([]string, error) {
	return []string{}, fmt.Errorf("This Guest driver dose not implement GetAttachDiskStatus")
}

func (self *SBaseGuestDriver) GetRebuildRootStatus() ([]string, error) {
	return []string{}, fmt.Errorf("This Guest driver dose not implement GetRebuildRootStatus")
}

func (self *SBaseGuestDriver) IsRebuildRootSupportChangeImage() bool {
	return true
}

func (self *SBaseGuestDriver) IsRebuildRootSupportChangeUEFI() bool {
	return true
}

func (self *SBaseGuestDriver) GetChangeConfigStatus(guest *models.SGuest) ([]string, error) {
	return []string{}, fmt.Errorf("This Guest driver dose not implement GetChangeConfigStatus")
}

func (self *SBaseGuestDriver) ValidateChangeConfig(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, cpuChanged bool, memChanged bool, newDisks []*api.DiskConfig) error {
	return nil
}

func (self *SBaseGuestDriver) ValidateDetachDisk(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, disk *models.SDisk) error {
	return nil
}

func (self *SBaseGuestDriver) ValidateCreateEip(ctx context.Context, userCred mcclient.TokenCredential, data jsonutils.JSONObject) error {
	return httperrors.NewInputParameterError("Not Implement ValidateCreateEip")
}

func (self *SBaseGuestDriver) ValidateResizeDisk(guest *models.SGuest, disk *models.SDisk, storage *models.SStorage) error {
	return fmt.Errorf("This Guest driver dose not implement ValidateResizeDisk")
}

func (self *SBaseGuestDriver) ValidateUpdateData(ctx context.Context, userCred mcclient.TokenCredential, data *jsonutils.JSONDict) (*jsonutils.JSONDict, error) {
	return data, nil
}

func (self *SBaseGuestDriver) GetDeployStatus() ([]string, error) {
	return []string{}, fmt.Errorf("This Guest driver dose not implement GetDeployStatus")
}

func (self *SBaseGuestDriver) IsNeedRestartForResetLoginInfo() bool {
	return true
}

func (self *SBaseGuestDriver) RequestDeleteDetachedDisk(ctx context.Context, disk *models.SDisk, task taskman.ITask, isPurge bool) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RqeuestSuspendOnHost(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RqeuestResumeOnHost(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) StartGuestResetTask(guest *models.SGuest, ctx context.Context, userCred mcclient.TokenCredential, isHard bool, parentTaskId string) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) StartGuestRestartTask(guest *models.SGuest, ctx context.Context, userCred mcclient.TokenCredential, isForce bool, parentTaskId string) error {
	data := jsonutils.NewDict()
	data.Set("is_force", jsonutils.NewBool(isForce))
	if err := guest.SetStatus(userCred, api.VM_STOPPING, ""); err != nil {
		return err
	}
	task, err := taskman.TaskManager.NewTask(ctx, "GuestRestartTask", guest, userCred, nil, parentTaskId, "", nil)
	if err != nil {
		return err
	}
	task.ScheduleRun(nil)
	return nil
}

func (self *SBaseGuestDriver) RequestSoftReset(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) AllowReconfigGuest() bool {
	return true
}

func (self *SBaseGuestDriver) DoGuestCreateDisksTask(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RequestChangeVmConfig(ctx context.Context, guest *models.SGuest, task taskman.ITask, instanceType string, vcpuCount, vmemSize int64) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) NeedRequestGuestHotAddIso(ctx context.Context, guest *models.SGuest) bool {
	return false
}

func (self *SBaseGuestDriver) RequestGuestHotAddIso(ctx context.Context, guest *models.SGuest, path string, boot bool, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RequestGuestHotRemoveIso(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RequestRebuildRootDisk(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RequestDiskSnapshot(ctx context.Context, guest *models.SGuest, task taskman.ITask, snapshotId, diskId string) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RequestDeleteSnapshot(ctx context.Context, guest *models.SGuest, task taskman.ITask, params *jsonutils.JSONDict) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RequestReloadDiskSnapshot(ctx context.Context, guest *models.SGuest, task taskman.ITask, params *jsonutils.JSONDict) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) RequestSyncToBackup(ctx context.Context, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement")
}

func (self *SBaseGuestDriver) GetMaxSecurityGroupCount() int {
	return 5
}

func (self *SBaseGuestDriver) getTaskRequestHeader(task taskman.ITask) http.Header {
	return task.GetTaskRequestHeader()
}

func (self *SBaseGuestDriver) IsSupportedBillingCycle(bc billing.SBillingCycle) bool {
	return true
}

func (self *SBaseGuestDriver) IsSupportPostpaidExpire() bool {
	return true
}

func (self *SBaseGuestDriver) RequestRenewInstance(guest *models.SGuest, bc billing.SBillingCycle) (time.Time, error) {
	return time.Time{}, nil
}

func (self *SBaseGuestDriver) IsSupportEip() bool {
	return false
}

func (self *SBaseGuestDriver) IsSupportPublicIp() bool {
	return false
}

func (self *SBaseGuestDriver) NeedStopForChangeSpec(guest *models.SGuest, cpuChanged, memChanged bool) bool {
	return false
}

func (self *SBaseGuestDriver) RemoteDeployGuestForCreate(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, host *models.SHost, desc cloudprovider.SManagedVMCreateConfig) (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (self *SBaseGuestDriver) RemoteActionAfterGuestCreated(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, host *models.SHost, ivm cloudprovider.ICloudVM, desc *cloudprovider.SManagedVMCreateConfig) {
	return
}

func (self *SBaseGuestDriver) RemoteDeployGuestForDeploy(ctx context.Context, guest *models.SGuest, ihost cloudprovider.ICloudHost, task taskman.ITask, desc cloudprovider.SManagedVMCreateConfig) (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (self *SBaseGuestDriver) RemoteDeployGuestForRebuildRoot(ctx context.Context, guest *models.SGuest, ihost cloudprovider.ICloudHost, task taskman.ITask, desc cloudprovider.SManagedVMCreateConfig) (jsonutils.JSONObject, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (self *SBaseGuestDriver) GetGuestInitialStateAfterCreate() string {
	return api.VM_READY
}

func (self *SBaseGuestDriver) GetGuestInitialStateAfterRebuild() string {
	return api.VM_READY
}

func (self *SBaseGuestDriver) IsNeedInjectPasswordByCloudInit(desc *cloudprovider.SManagedVMCreateConfig) bool {
	return false
}

func (self *SBaseGuestDriver) GetWindowsUserDataType() string {
	return cloudprovider.CLOUD_POWER_SHELL
}

func (self *SBaseGuestDriver) IsWindowsUserDataTypeNeedEncode() bool {
	return false
}

func (self *SBaseGuestDriver) IsSupportdDcryptPasswordFromSecretKey() bool {
	return true
}

func (self *SBaseGuestDriver) GetUserDataType() string {
	return cloudprovider.CLOUD_CONFIG
}

func (self *SBaseGuestDriver) GetDefaultAccount(desc cloudprovider.SManagedVMCreateConfig) string {
	if strings.ToLower(desc.OsType) == strings.ToLower(osprofile.OS_TYPE_WINDOWS) {
		return api.VM_DEFAULT_WINDOWS_LOGIN_USER
	}
	return api.VM_DEFAULT_LINUX_LOGIN_USER
}

func (self *SBaseGuestDriver) OnGuestChangeCpuMemFailed(ctx context.Context, guest *models.SGuest, data *jsonutils.JSONDict, task taskman.ITask) error {
	return nil
}

func (self *SBaseGuestDriver) RequestSyncConfigOnHost(ctx context.Context, guest *models.SGuest, host *models.SHost, task taskman.ITask) error {
	return fmt.Errorf("SBaseGuestDriver: Not Implement")
}

func (self *SBaseGuestDriver) IsSupportGuestClone() bool {
	return true
}

func (self *SBaseGuestDriver) RequestSyncSecgroupsOnHost(ctx context.Context, guest *models.SGuest, host *models.SHost, task taskman.ITask) error {
	return nil // do nothing
}

func (self *SBaseGuestDriver) CancelExpireTime(
	ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest) error {
	return guest.CancelExpireTime(ctx, userCred)
}

func (self *SBaseGuestDriver) IsSupportPublicipToEip() bool {
	return false
}

func (self *SBaseGuestDriver) RequestConvertPublicipToEip(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestConvertPublicipToEip")
}

func (self *SBaseGuestDriver) IsSupportSetAutoRenew() bool {
	return false
}

func (self *SBaseGuestDriver) RequestSetAutoRenewInstance(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, autoRenew bool, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestSetAutoRenewInstance")
}

func (self *SBaseGuestDriver) IsSupportMigrate() bool {
	return false
}

func (self *SBaseGuestDriver) IsSupportLiveMigrate() bool {
	return false
}

func (self *SBaseGuestDriver) CheckMigrate(guest *models.SGuest, userCred mcclient.TokenCredential, input api.GuestMigrateInput) error {
	return httperrors.NewNotAcceptableError("Not allow for hypervisor %s", guest.GetHypervisor())
}

func (self *SBaseGuestDriver) CheckLiveMigrate(guest *models.SGuest, userCred mcclient.TokenCredential, input api.GuestLiveMigrateInput) error {
	return httperrors.NewNotAcceptableError("Not allow for hypervisor %s", guest.GetHypervisor())
}
func (self *SBaseGuestDriver) RequestMigrate(ctx context.Context, guest *models.SGuest, userCred mcclient.TokenCredential, data *jsonutils.JSONDict, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestMigrate")
}

func (self *SBaseGuestDriver) RequestLiveMigrate(ctx context.Context, guest *models.SGuest, userCred mcclient.TokenCredential, data *jsonutils.JSONDict, task taskman.ITask) error {
	return fmt.Errorf("Not Implement RequestLiveMigrate")
}

func (self *SBaseGuestDriver) RequestRemoteUpdate(ctx context.Context, guest *models.SGuest, userCred mcclient.TokenCredential, replaceTags bool) error {
	// nil ops
	return nil
}

func (self *SBaseGuestDriver) ValidateRebuildRoot(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest, input *api.ServerRebuildRootInput) (*api.ServerRebuildRootInput, error) {
	return input, nil
}

func (self *SBaseGuestDriver) ValidateDetachNetwork(ctx context.Context, userCred mcclient.TokenCredential, guest *models.SGuest) error {
	return nil
}
