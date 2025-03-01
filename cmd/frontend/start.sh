#!/bin/sh
#@auth cl
#@time 20240920

# 获取运行服务器架构名称
ARCH=$(uname -m)
# 获取操作系统
OS=$(sw_vers -productName)
if [ -z ${OS} ]; then
    OS="unknown"
fi
TARGET=frontend

#echo -e "\033[30m ### 30:黑   ### \033[0m"
#echo -e "\033[31m ### 31:红   ### \033[0m"ls


#echo -e "\033[32m ### 32:绿   ### \033[0m"
#echo -e "\033[33m ### 33:黄   ### \033[0m"
#echo -e "\033[34m ### 34:蓝色 ### \033[0m"
#echo -e "\033[35m ### 35:紫色 ### \033[0m"
#echo -e "\033[36m ### 36:深绿 ### \033[0m"
#echo -e "\033[37m ### 37:白色 ### \033[0m"

# 获取shell脚本运行路径
SHELL_BASE_PATH=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)

# 构建包路径
BULID_PKG_PATH=${SHELL_BASE_PATH}/build

# 时间
DATE=$(date "+%Y%m%d%H%M%S")

# 检查结果函数
function check_result()
{
    local result=$1
    local object=$2
    if [ ${result} -ne 0 ]; then
        echo -e "\033[31m ### [check_result] result is failed! result:${result}, msg:${object}   ### \033[0m"
        exit 1
    fi
    echo -e "\033[32m ### [check_result] result is success! result:${result}, msg:${object}   ### \033[0m"
}

# 检查是否为文件函数
function check_file()
{
    local file_path=$1
    if [ ! -f ${file_path} ] ; then
        echo -e "\033[31m ### [check_file] ${file_path} is not exist!   ### \033[0m"
        print_help
        exit 1
    fi
}

# 检查是否为目录函数
function check_dir() {
    local file_path=$1
    if [ ! -d ${file_path} ] ; then
        echo -e "\033[31m ### [check_dir] ${file_path} is not exist!   ### \033[0m"
        print_help
        exit 1
    fi
}

function resources() {
    echo -e "\033[33m ### [resources]   ### \033[0m"
    local templ_dir=${SHELL_BASE_PATH}/internal/gin/templates
    if [ ! -d ${templ_dir} ] ; then
      echo -e "\033[31m ### ${templ_dir} is not exist!   ### \033[0m"
      return
    fi
    cp -rvf ${templ_dir} ${BULID_PKG_PATH}
    echo -e "\033[32m ### [2-1] copy ${templ_dir} to ${BULID_PKG_PATH}   ### \033[0m"

    local static_dir=${SHELL_BASE_PATH}/internal/gin/static
    if [ ! -d ${static_dir} ] ; then
      echo -e "\033[31m ### ${static_dir} is not exist!   ### \033[0m"
      return
    fi
    cp -rvf ${static_dir} ${BULID_PKG_PATH}
    echo -e "\033[32m ### [2-2] copy ${static_dir} to ${BULID_PKG_PATH}   ### \033[0m"
}

function build() {
    clean
    echo -e "\033[33m ### [build]   ### \033[0m"
    rm -rf ${BULID_PKG_PATH}
    mkdir -p ${BULID_PKG_PATH}
    go build -o ${TARGET}
    check_file ${TARGET}

    mv ${TARGET} ${BULID_PKG_PATH}
    echo -e "\033[32m ### [2-1] move ${TARGET} to ${BULID_PKG_PATH}   ### \033[0m"

    cp -rvf ${SHELL_BASE_PATH}/conf ${BULID_PKG_PATH}
    echo -e "\033[32m ### [2-2] copy ${SHELL_BASE_PATH}/conf to ${BULID_PKG_PATH}   ### \033[0m"

    if [ ${TARGET} = "frontend" ] ; then
      resources
    fi
}

function run() {
    echo -e "\033[33m ### [run]   ### \033[0m"
    check_dir ${BULID_PKG_PATH}
    cd ${BULID_PKG_PATH}
    check_file ${TARGET}
    chmod +x ${TARGET}
    echo -e "\033[32m ### [1-1] ./${TARGET} server   ### \033[0m"
    ./${TARGET} server
}

function clean() {
    echo -e "\033[33m ### [clean]   ### \033[0m"
    rm -rf ${TARGET}
    echo -e "\033[32m ### [2-1] remove ${TARGET}   ### \033[0m"
    rm -rf ${BULID_PKG_PATH}
    echo -e "\033[32m ### [2-2] remove ${BULID_PKG_PATH}   ### \033[0m"

}

function copy_script() {
    local dest_dir=$1
    local script_name=run.sh
    echo -e "\033[33m ### [copy_script]   ### \033[0m"
    pushd ../../script
    check_file ${script_name}
    chmod +x ${script_name}
    cp -rvf ${script_name} ${dest_dir}
    check_file ${dest_dir}/${script_name}
    echo -e "\033[32m ### copy script/run.sh to ${dest_dir}   ### \033[0m"
    popd

    pushd ${dest_dir}
    # 使用 sed 命令进行替换
    local old_target="TARGET=xreplacex"
    local new_target="TARGET=${TARGET}"
    check_file ${script_name}
    if [ ${OS} = "macOS" ]; then
        sed -i '' "s/${old_target}/${new_target}/g" ${script_name}
    else
        sed -i "s/${old_target}/${new_target}/g" ${script_name}
    fi
    check_result $? "sed ${old_target} to ${new_target} in ${script_name}, OS=${OS}"
    echo -e "\033[32m ### replace ${old_target} to ${new_target} in ${dest_dir}/${script_name} ### \033[0m"
    popd
}

function compress() {
    local pkg_name=$1
    echo -e "\033[33m ### [compress]   ### \033[0m"
    check_dir ${TARGET}
    tar -zcvf ${pkg_name} ${TARGET}
    check_file ${SHELL_BASE_PATH}/${pkg_name}
}

function package() {
    echo -e "\033[33m ### [package]   ### \033[0m"
    check_dir ${BULID_PKG_PATH}
    check_file ${BULID_PKG_PATH}/${TARGET}

    # 1-1、构建临时打包目录
    rm -rf ${SHELL_BASE_PATH}/${TARGET}
    mkdir -p ${SHELL_BASE_PATH}/${TARGET}
    echo -e "\033[32m ### [5-1-1] mkdir -p ${SHELL_BASE_PATH}/${TARGET}   ### \033[0m"

    # 2-1、复制运行文件到临时打包目录
    cp -rvf ${BULID_PKG_PATH}/* ${SHELL_BASE_PATH}/${TARGET}
    echo -e "\033[32m ### [5-2-2] copy ${BULID_PKG_PATH}/* to ${SHELL_BASE_PATH}/${TARGET}   ### \033[0m"

    # 3-1、第一次打包
    local pkg1_name=${TARGET}-${DATE}.tar.gz
    tar -zcvf ${SHELL_BASE_PATH}/${pkg1_name} ${TARGET}
    check_file ${SHELL_BASE_PATH}/${pkg1_name}
    echo -e "\033[32m ### [5-3-1] package ${SHELL_BASE_PATH}/${pkg1_name}   ### \033[0m"

    # 3-2、创建打包目录
    rm -rf ${BULID_PKG_PATH}/package
    mkdir -p ${BULID_PKG_PATH}/package
    echo -e "\033[32m ### [5-3-2] reset ${BULID_PKG_PATH}/package   ### \033[0m"

    # 3-3、移动压缩文件到打包目录
    mv ${SHELL_BASE_PATH}/${pkg1_name} ${BULID_PKG_PATH}/package
    echo -e "\033[32m ### [5-3-3] move ${SHELL_BASE_PATH}/${pkg1_name} to ${BULID_PKG_PATH}/package   ### \033[0m"

    # 3-4、删除临时打包目录
    rm -rf ${SHELL_BASE_PATH}/${TARGET}
    echo -e "\033[32m ### [5-3-4] remove ${SHELL_BASE_PATH}/${TARGET}   ### \033[0m"

    # 4-1、复制script文件到打包目录
    copy_script ${BULID_PKG_PATH}/package
    echo -e "\033[32m ### [5-4-1] copy script/run.sh to ${BULID_PKG_PATH}/package   ### \033[0m"

    # 5-1、第二次打包
    local pkg2_name=package-${DATE}.tar.gz
    pushd ${BULID_PKG_PATH}
    tar -zcvf ${BULID_PKG_PATH}/${pkg2_name} package
    popd
    check_file ${BULID_PKG_PATH}/${pkg2_name}
    echo -e "\033[32m ### [5-5-1] package ${BULID_PKG_PATH}/${pkg2_name}   ### \033[0m"
}

function print_help() {
    echo -e "\033[35m ######################### HELP ARCH:${ARCH} OS:${OS} ######################### \033[0m"
    echo -e "\033[35m #sh start.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m        -b       : build \033[0m"
    echo -e "\033[35m        -r       : run  \033[0m"
    echo -e "\033[35m        -c       : clean \033[0m"
    echo -e "\033[35m        -p       : package \033[0m"
    echo -e "\033[35m        -        : build ->package -> run \033[0m"
    echo -e "\033[35m        - help   : help \033[0m"
    echo -e "\033[35m ######################### HELP ARCH:${ARCH} OS:${OS} ######################### \033[0m"
    exit 1
}

function main() {
    echo -e "\033[34m ######################### ${TARGET} input param is $@ ######################### \033[0m"
    case $1 in
        "-b")
            build
            ;;
        "-r")
            run
            ;;
        "-c")
            clean
            ;;
        "-p")
            package
            ;;
        "-")
            build
            package
            run
          ;;
        *)
          print_help
          ;;
    esac
}

main "$@"