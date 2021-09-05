# 当前tag版本
version=`git describe --tags $(git rev-list --tags --max-count=1 --branches master)`
binary="lmap"
echo "build version: $version"

# cross_compiles
make all

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        binary_dir_name="${binary}_${version}_${os}_${arch}" #压缩包目录
        binary_path="./packages/${binary_dir_name}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./${binary}_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${binary_path}
            mv ./${binary}_${os}_${arch}.exe ${binary_path}/${binary}.exe
        else
            if [ ! -f "./${binary}_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${binary_path}
            mv ${binary}_${os}_${arch} ${binary_path}/${binary}
        fi  
        cp ../LICENSE ${binary_path}
        cp ../README.md ${binary_path}

        # packages
        cd ./packages
        tar -zcf ${binary_dir_name}.tar.gz ${binary_dir_name}
        cd ..
        rm -rf ${binary_path}
    done
done

cd -
