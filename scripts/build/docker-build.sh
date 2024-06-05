#!/usr/bin/bash

# 重新vendor
refresh_vendor(){
    mkdir -p scripts/build/vendor
    cp go.mod scripts/build/
    go mod vendor
    # https://deepinout.com/linux/linux-ask-and-questions/153_tk_1704111721.html  无法覆盖mv，只能先cp 在rm
    cp -rf vendor/* scripts/build/vendor/
    rm -rf vendor
}

start=$(date "+%H:%M:%S")
cd ../..

VERSION=$(date "+%Y%m%d") #一个版本变一次,# $(date "+%Y-%m-%d %H:%M:%S")
COMMIT=$(git rev-parse --verify HEAD)
NOW=$(date '+%FT%T%z')

# -f 参数判断 $file 是否存在
if [ ! -f "scripts/build/go.mod" ]; then
    echo 旧go.mod不存在
    refresh_vendor
else
    echo 旧go.mod存在
    # 根据go mod变化，更新vendor
    modMd5New=`md5sum go.mod|awk '{print $1}'`
    modMd5Old=`md5sum scripts/build/go.mod|awk '{print $1}'`
    if [ "$modMd5New" != "$modMd5Old" ]; then
        echo "旧go.mod已改变" \| new $modMd5New \| old $modMd5Old
        refresh_vendor
    else
        echo "旧go.mod未改变" $modMd5New
    fi
fi

mkdir -p vendor/
cp -r scripts/build/vendor/* vendor/ \
&& docker build --tag sonic:${VERSION} --file scripts/Dockerfile \
--build-arg BUILD_COMMIT=${COMMIT} --build-arg SONIC_VERSION=${VERSION} --build-arg BUILD_TIME=${NOW}  \
.  && rm -rf vendor

echo $(date "+%Y-%m-%d") build耗时 $start '~' $(date "+%H:%M:%S") \| VERSION/tag[$VERSION] \| COMMIT[$COMMIT] \| NOW[$NOW]