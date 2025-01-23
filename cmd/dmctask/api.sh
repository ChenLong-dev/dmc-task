#!/bin/sh
#@auth cl
#@time 20240920

SHELL_BASE_PATH=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)

TASK_API_FILE=${SHELL_BASE_PATH}/task.api

TYPES_PATH=${SHELL_BASE_PATH}/api/internal/types/types.go

DEST_API_FILE=${SHELL_BASE_PATH}/../../core/common/common.go

DEST_API_DIR=${SHELL_BASE_PATH}/../../core/common


#echo -e "\033[30m ### 30:黑   ### \033[0m"
#echo -e "\033[31m ### 31:红   ### \033[0m"
#echo -e "\033[32m ### 32:绿   ### \033[0m"
#echo -e "\033[33m ### 33:黄   ### \033[0m"
#echo -e "\033[34m ### 34:蓝色 ### \033[0m"
#echo -e "\033[35m ### 35:紫色 ### \033[0m"
#echo -e "\033[36m ### 36:深绿 ### \033[0m"
#echo -e "\033[37m ### 37:白色 ### \033[0m"


# 检查是否为文件函数
function check_file()
{
    local file_path=$1
    if [ ! -f ${file_path} ] ; then
        echo -e "\033[31m ### [check_file] ${file_path} is not exist!   ### \033[0m"
        exit 1
    fi
}

function generate_api() {
    echo ${SHELL_BASE_PATH}

    check_file ${TASK_API_FILE}

    goctl api go -api task.api --dir ./api --style=gozero --home ../../goctl

    check_file ${TYPES_PATH}

    cp -rvf ${TYPES_PATH} ${DEST_API_FILE}

    pushd ${DEST_API_DIR}
    sed -i "s/package types/package common/g" common.go  # windows
    # sed -i '' "s/package types/package common/g" common.go  # linux

    echo -e "\033[32m ### -------------------------------------------------   ### \033[0m"
    head -5 common.go
    echo -e "\033[32m ### -------------------------------------------------   ### \033[0m"
    popd
}

function main() {
    generate_api
}

main "$@"
