GREEN='\033[0;32m'
RED='\033[0;31m'
NOCOLOR='\033[0m'

sh build.sh && echo "${GREEN}BUILD SUCCESS!${NOCOLOR}"

SUCCESS_INFO="EXECUTE STMTS SUCCESS!"
for i in ${PWD}/lox_demo/*; do
    name=`echo "$i"|rev|cut -d "/" -f1|rev`
    # echo $name
    echo "${GREEN}case ${name}${NOCOLOR}, result: \c"  && ./main lox_demo/${name} | grep -q "${SUCCESS_INFO}" && echo "${GREEN}success${NOCOLOR}" || echo "${RED}failed${NOCOLOR}"
done