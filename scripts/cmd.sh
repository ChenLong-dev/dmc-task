#!/bin/sh
#@auth cl
#@time 20240920

# 获取运行服务器架构名称
ARCH=$(uname -m)

#echo -e "\033[30m ### 30:黑   ### \033[0m"
#echo -e "\033[31m ### 31:红   ### \033[0m"
#echo -e "\033[32m ### 32:绿   ### \033[0m"
#echo -e "\033[33m ### 33:黄   ### \033[0m"
#echo -e "\033[34m ### 34:蓝色 ### \033[0m"
#echo -e "\033[35m ### 35:紫色 ### \033[0m"
#echo -e "\033[36m ### 36:深绿 ### \033[0m"
#echo -e "\033[37m ### 37:白色 ### \033[0m"

# 获取shell脚本运行路径
# shellcheck disable=SC2034
# shellcheck disable=SC2039
SHELL_BASE_PATH=$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd)


# shellcheck disable=SC2112
function print_help() {
    echo -e "\033[35m ######################### HELP ARCH:${ARCH} ######################### \033[0m"
    echo -e "\033[35m #sh scp_file.sh {param} \033[0m"
    echo -e "\033[35m {param}: \033[0m"
    echo -e "\033[35m        -b       : build \033[0m"
    echo -e "\033[35m        -r       : run  \033[0m"
    echo -e "\033[35m        -c       : clean \033[0m"
    echo -e "\033[35m        -p       : package \033[0m"
    echo -e "\033[35m        -        : build -> package \033[0m"
    echo -e "\033[35m        - help   : help \033[0m"
    echo -e "\033[35m ######################### HELP ARCH:${ARCH} ######################### \033[0m"
    exit 1
}

# shellcheck disable=SC2112
function test() {
    echo -e "\033[32m test sleep 5s \033[0m"
    sleep 5
    echo -e "\033[32m test end \033[0m"
}


# shellcheck disable=SC2112
function main() {
    echo -e "\033[34m ######################### build.sh input param is $@ ######################### \033[0m"
    case $1 in
        "-b")
          echo -e "\033[34m build ${@:2} \033[0m"
            ;;
        "-h")
          print_help
            ;;
        *)
          print_help
          ;;
    esac
}

main "$@"