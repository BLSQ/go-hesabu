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

for name in test/bad_*.json
do
    if bin/hesabucli $name >/dev/null 2>&1
    then
        $its_all_good=1
        echo "${name} \033[1;31mFAIL\033[0m"
    else
        echo "${name} \033[1;32mPASS\033[0m"
    fi
done

for name in $(ls -1 test/*.json | grep -v 'bad_')
do
    if bin/hesabucli $name >/dev/null 2>&1
    then
        echo "${name} \033[1;32mPASS\033[0m"
    else
        $its_all_good=1
        echo "${name} \033[1;31mFAIL\033[0m"
    fi
done

if [ "${its_all_good}" -gt "0" ]
then
    echo "Some examples are failing"
fi

exit $its_all_good