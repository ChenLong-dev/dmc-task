#!/bin/sh
#@auth cl
#@time 20240920

# 获取运行服务器架构名称
ARCH=$(uname -m)

SHELL_BASE_PATH=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)

TASK_API_FILE=${SHELL_BASE_PATH}/biz.api

#echo -e "\033[30m ### 30:黑   ### \033[0m"
#echo -e "\033[31m ### 31:红   ### \033[0m"
#echo -e "\033[32m ### 32:绿   ### \033[0m"
#echo -e "\033[33m ### 33:黄   ### \033[0m"
#echo -e "\033[34m ### 34:蓝色 ### \033[0m"
#echo -e "\033[35m ### 35:紫色 ### \033[0m"
#echo -e "\033[36m ### 36:深绿 ### \033[0m"
#echo -e "\033[37m ### 37:白色 ### \033[0m"

# 检查是否为文件函数
# shellcheck disable=SC2112
function check_file()
{
    local file_path=$1
    if [ ! -f ${file_path} ] ; then
        echo -e "\033[31m ### [check_file] ${file_path} is not exist!   ### \033[0m"
        exit 1
    fi
}

# shellcheck disable=SC2112
function generate_api() {
    echo ${SHELL_BASE_PATH}

    check_file ${TASK_API_FILE}

    goctl api go -api biz.api --dir ./biz --style=gozero --home ../../goctl
}

# shellcheck disable=SC2112
function print_help() {
    echo -e "\033[35m ######################### HELP ARCH:${ARCH} ######################### \033[0m"
    echo -e "\033[35m #sh api.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m        -api      \033[0m"
    echo -e "\033[35m        -help     \033[0m"
    echo -e "\033[35m ######################### HELP ARCH:${ARCH} ######################### \033[0m"
    exit 1
}

# shellcheck disable=SC2112
function main() {
  echo -e "\033[34m ######################### api.sh input param is $@ ######################### \033[0m"
      case $1 in
          "-api")
            generate_api
              ;;
          "-help")
            print_help
              ;;
          *)
            print_help
            ;;
      esac
}

main "$@"
