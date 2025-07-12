# -------- 在本地的环境进行 ---------
# 下载交叉编译所需要的依赖
brew install FiloSottile/musl-cross/musl-cross
# 切换到脚本目录

# 执行编译打包
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build .

# 重命名
mv XXX $your_name$

# 压缩，拆分成多个压缩包，上传到pod更快
tar cjf - $your_name$ |split -b 5m - $your_name$.tar.bz2.




# -------- 在pod中 ---------
# 安装lrzsz，用于文件的上传和下载
apt install -y lrzsz

# 安装解压文件
apt install -y bzip2

# 上选自己需要的文件，选择目录
rz

# 解压缩：
cat $your_name$.tar.bz2.* | tar xj

# 给文件赋执行权限
chmod +x $your_name$