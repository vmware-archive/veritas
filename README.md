# Veritas

Veritas is a cli for getting at Diego's truth.

For a linux build on a bosh vm:

```bash
pushd $HOME
rm veritas

wget http://onsi-public.s3.amazonaws.com/veritas
chmod +x ./veritas
export PATH=$PATH:$PWD

veritas autodetect && `veritas autodetect`

veritas completions > ./.veritas_completions.bash
source ./.veritas_completions.bash
popd
```
