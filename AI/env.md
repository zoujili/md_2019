#Nvidia driver
 ##uninstall old
sudo apt-get remove –purge nvidia*
 ##stop nouveau
sudo gedit /etc/modprobe.d/blacklist-nouveau.conf  
Add:
```
blacklist nouveau
option nouveau modeset=0
```
sudo update-initramfs -u  
lsmod | grep nouveau  
```
 should none
```
 ##install driver 
 ```
  231  sudo add-apt-repository ppa:graphics-drivers/ppa
  232  sudo apt-get install --reinstall python3-apt
  233  sudo add-apt-repository ppa:graphics-drivers/ppa
  234  sudo apt-get remove --purge python-apt
  235  sudo apt-get install python-apt -f
  236  sudo find / -name "apt_pkg.cpython-35m-x86_64-linux-gnu.so"
  237  cd /usr/lib/python3/dist-packages/
  238  sudo cp apt_pkg.cpython-35m-x86_64-linux-gnu.so apt_pkg.cpython-36m-x86_64-linux-gnu.so 
  239  sudo add-apt-repository ppa:graphics-drivers/ppa
  240  sudo apt-get update
  241  sudo apt-get install nvidia-410
  242  sudo apt-get update 
  243  sudo apt-get upgrade

```
nvidia-smi

 sudo reboot
 F9 ->  secure boot -> disabled  
#Cuda

https://developer.nvidia.com/cuda-toolkit-archive  
sudo apt-get install freeglut3-dev build-essential libx11-dev libxmu-dev libglu1  
sudo sh cuda_9.0.176_384.81_linux.run  
sudo vim ~/.bashrc  # 配置环境变量  

文末追加以下三行代码  
export PATH=/usr/local/cuda-9.0/bin:$PATH  
export LD_LIBRARY_PATH=/usr/local/cuda-9.0/lib64:$LD_LIBRARY_PATH
export CUDA_HOME=/usr/local/cuda

source ~/.bashrc  
sudo ldconfig
nvcc --version

cd /usr/local/cuda-9.0/samples
root用户下执行make -j
cd ./bin/x86_64/linux/release
./deviceQuery

#CudaNN

https://developer.nvidia.com/rdp/cudnn-archive
cd Download
tar -xzvf cudnn-9.0-linux-x64-v7.tgz
sudo cp cuda/include/cudnn.h /usr/local/cuda/include
sudo cp cuda/lib64/libcudnn* /usr/local/cuda/lib64
sudo chmod a+r /usr/local/cuda/include/cudnn.h /usr/local/cuda/lib64/libcudnn*
sudo dpkg -i libcudnn7_7.0.3.11-1+cuda9.0_amd64.deb
sudo dpkg -i libcudnn7-devel_7.0.3.11-1+cuda9.0_amd64.deb
sudo dpkg -i libcudnn7-doc_7.0.3.11-1+cuda9.0_amd64.deb

test:
cp -r /usr/src/cudnn_samples_v7/ $HOME
cd $HOME/cudnn_samples_v7/mnistCUDNN
make clean && make
./mnistCUDNN


#Pip3 update
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python get-pip.py

#Tensorflow
pip3 install --upgrade tensorflow-gpu==1.8.0


#Pycharm
maybe can't find cuda 9.0 ->LIBRARY_PATH=/usr/local/cuda-9.0/lib64

#Keras
pip install keras==2.1.1


#Check GPU
watch -n 2  nvidia-smi