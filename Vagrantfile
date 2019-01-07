# Recommended:
# 	vagrant plugin install vagrant-vbguest

MEMORY_MB = 4096
CPU_CORES = 2

Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"
  config.vm.provider :virtualbox do |virtualbox, override|
    virtualbox.memory = MEMORY_MB
    virtualbox.cpus = CPU_CORES
    virtualbox.gui = false
    override.vm.box_download_checksum_type = "sha256"
    override.vm.box_download_checksum = "b24c912b136d2aa9b7b94fc2689b2001c8d04280cf25983123e45b6a52693fb3"
    override.vm.box_url = "https://cloud.centos.org/centos/7/vagrant/x86_64/images/CentOS-7-x86_64-Vagrant-1803_01.VirtualBox.box"
  end
  config.vm.provision "ansible" do |ansible|
    ansible.playbook = "vagrant-playbook.yml"
    ansible.verbose = true
  end

  config.vm.synced_folder ".", "/vagrant", type: "rsync",
    rsync__exclude: [
      "bazel-bazel-imagick/",
      "bazel-out/",
      "bazel-bin/",
      "bazel-genfiles/",
      "bazel-testlogs/",
    ]

  config.vm.network "private_network", type: "dhcp"

  # Mitigate "VBoxHeadless + logd using all available CPU"
  $enable_serial_logging = false
end
