function VolumeViewModel(data) {
  this.Id = data.Id;
  this.Name = data.Name;
  this.Driver = data.Driver;
  this.Mountpoint = data.Mountpoint;
  if (data.Whale) {
    this.Metadata = {};
    if (data.Whale.ResourceControl) {
      this.Metadata.ResourceControl = {
        OwnerId: data.Whale.ResourceControl.OwnerId
      };
    }
  }
}
