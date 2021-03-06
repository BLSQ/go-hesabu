#! /bin/sh
#
#
set -e

# Loop over all json files in the test folder, feed them through the
# binary and check if hesabu exits the right way:
#
# - Everything starting with bad should exit with error
# - Everything else should exit without error.

its_all_good=0

cli=bin/hesabucli
if [[ "$OSTYPE" == "darwin"* ]]; then
    cli=bin/hesabucli-mac
fi

for name in test/bad_*.json
do
    if $cli $name >/dev/null 2>&1
    then
        its_all_good=1
        echo "${name} \033[1;31mFAIL\033[0m"
    else
        echo "${name} \033[1;32mPASS\033[0m"
    fi
done

for name in $(ls -1 test/*.json | grep -v 'bad_')
do
    if $cli $name >/dev/null 2>&1
    then
        echo "${name} \033[1;32mPASS\033[0m"
    else
        its_all_good=1
        echo "${name} \033[1;31mFAIL\033[0m"
    fi
done

if [ "${its_all_good}" -gt "0" ]
then
    echo "\n\033[1;31mSome examples are failing\033[0m"
else
    echo "\n\033[1;32mIt's all good\033[0m"
fi

exit $its_all_good
