import sys
import time
from os import system, chdir
from utils import execute_commands

OS = ['darwin', 'linux', 'windows']
ARCHS = ['386', 'amd64']

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('Usage: python3 build.py vx.y.z [b6k_cs_pw]')
        exit(-1)

    version = sys.argv[1]
    b6k_cs_pw = sys.argv[2] if len(sys.argv) > 2 else None
    versioned_dist = 'slangd-' + version.replace('.', '_')
    build_time = int(time.time())

    ldflags = f"-X main.Version={version} "
    ldflags += f"-X main.BuildTime={build_time} "

    for os in OS:
        for arch in ARCHS:
            filename_with_ending = filename = f"{versioned_dist}-{os}-{arch}"
            if os == 'windows':
                filename_with_ending += ".exe"
                compress_cmd = f"zip {filename}.zip {filename_with_ending}"
            else:
                compress_cmd = f"tar -czvf {filename}.tar.gz {filename_with_ending}"

            execute_commands([
                f"env GOOS={os} GOARCH={arch} go build -ldflags \"{ldflags}\" -o ./ci/release/{filename_with_ending} ./cmd/slangd",
            ])

            if os == 'windows' and b6k_cs_pw:
                execute_commands([
                    f"osslsigncode sign -pkcs12 ./ci/b6k_csc.p12 -pass {b6k_cs_pw} -in ./ci/release/{filename_with_ending} -out ./ci/release/signed_{filename_with_ending}",
                ], True, False)
                execute_commands([
                    f"rm ./ci/release/{filename_with_ending}",
                    f"mv ./ci/release/signed_{filename_with_ending} ./ci/release/{filename_with_ending}",
                ])

            chdir("./ci/release/")
            execute_commands([
                compress_cmd,
                f"rm {filename_with_ending}",
            ])
            chdir("../..")
