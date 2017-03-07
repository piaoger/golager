

OS=$(uname)
APP_NAME=lager
APP_BIN_DIR=$(cd ${0%/*} && echo $PWD)/..
APP_ROOT=$APP_BIN_DIR/..

APP_CACHE_DIR=~/.$APP_NAME/cache

mkdir -p $APP_CACHE_DIR
cd $APP_CACHE_DIR


# -----------------------------------------
# golang
# -----------------------------------------
echo "installing golang"
rm -rf $APP_CACHE_DIR/golang

GOLANG_VERSION=1.8
if [ "$OS" = "Darwin" ]; then
   export GOLANG_PKG_NAME=go$GOLANG_VERSION.darwin-amd64.tar.gz
else
   export GOLANG_PKG_NAME=go$GOLANG_VERSION.linux-amd64.tar.gz
fi

GOLANG_DOWNLOAD_URL=http://olxkpcnfn.bkt.clouddn.com/$GOLANG_PKG_NAME
#GOLANG_DOWNLOAD_URL=https://storage.googleapis.com/golang/$GOLANG_PKG_NAME

if [ ! -f $GOLANG_PKG_NAME ]; then
    wget $GOLANG_DOWNLOAD_URL
fi

tar -xf $GOLANG_PKG_NAME

APP_BIN_GO_DIR=$APP_BIN_DIR/apps/golang

rm -rf $APP_BIN_GO_DIR
mkdir -p $APP_BIN_GO_DIR

mv ./go/bin $APP_BIN_GO_DIR/$OS
mv ./go/pkg $APP_BIN_GO_DIR/pkg
mv ./go/src $APP_BIN_GO_DIR/src
mv ./go/VERSION $APP_BIN_GO_DIR/VERSION

# remove decompressed go folder
rm -rf $APP_CACHE_DIR/golang

# -----------------------------------------
# golang x/net
# -----------------------------------------


echo "installing golang.com/x/net"
cd $APP_CACHE_DIR

NET_RELEASE_PACKAGE=release-branch.go1.8.zip
NET_RELEASE_PACKAGE_NAME=net-release-branch.go1.8

NET_RELEASE_PACKAGE_DOWNLOAD_URL=https://github.com/golang/net/archive/$NET_RELEASE_PACKAGE
if [ ! -f $NET_RELEASE_PACKAGE_NAME.zip ]; then
    wget $NET_RELEASE_PACKAGE_DOWNLOAD_URL -O $NET_RELEASE_PACKAGE_NAME.zip
fi

rm -rf ${APP_ROOT}/vendor/src/golang.org/x/
mkdir -p ${APP_ROOT}/vendor/src/golang.org/x/
cp ${APP_CACHE_DIR}/$NET_RELEASE_PACKAGE_NAME.zip ${APP_ROOT}/vendor/src/golang.org/x/

cd ${APP_ROOT}/vendor/src/golang.org/x/
unzip -q $NET_RELEASE_PACKAGE_NAME.zip

mv ${APP_ROOT}/vendor/src/golang.org/x/$NET_RELEASE_PACKAGE_NAME ${APP_ROOT}/vendor/src/golang.org/x/net
rm $NET_RELEASE_PACKAGE_NAME.zip

# -----------------------------------------
# golang x/net
# -----------------------------------------


# echo "installing golang.com/x/text"
# cd $APP_CACHE_DIR

# NET_RELEASE_PACKAGE=master.zip
# NET_RELEASE_PACKAGE_NAME=golang-text-master

# NET_RELEASE_PACKAGE_DOWNLOAD_URL=https://github.com/golang/text/archive/$NET_RELEASE_PACKAGE
# if [ ! -f $NET_RELEASE_PACKAGE_NAME.zip ]; then
#     wget $NET_RELEASE_PACKAGE_DOWNLOAD_URL -O $NET_RELEASE_PACKAGE_NAME.zip
# fi

# rm -rf ${APP_ROOT}/vendor/src/golang.org/x/
# mkdir -p ${APP_ROOT}/vendor/src/golang.org/x/
# cp ${APP_CACHE_DIR}/$NET_RELEASE_PACKAGE_NAME.zip ${APP_ROOT}/vendor/src/golang.org/x/

# cd ${APP_ROOT}/vendor/src/golang.org/x/
# unzip -q $NET_RELEASE_PACKAGE_NAME.zip

# mv ${APP_ROOT}/vendor/src/golang.org/x/$NET_RELEASE_PACKAGE_NAME ${APP_ROOT}/vendor/src/golang.org/x/text
# rm $NET_RELEASE_PACKAGE_NAME.zip
