#!/usr/bin/python3

import argparse
import json
import os
import subprocess
import sys
from queue import Queue
from threading import Thread

def build(queue):
    while True:
        arch = queue.get()
        goos = arch["GOOS"]
        goarch = arch["GOARCH"]
        args = arch["args"]
        template = arch["template"]
        print(f"Building {goos}/{goarch}")
        path = template.replace("{{.Dir}}", os.path.basename(os.getcwd())).replace("{{.OS}}", goos).replace("{{.Arch}}",goarch)
        if goos == "windows":
            path += ".exe"
        cmd = ["go", "build", "-o", path] + args
        env = os.environ.copy()
        env["GOOS"] = goos
        env["GOARCH"] = goarch
        subprocess.run(cmd, env=env)
        queue.task_done()

def filter(archs, osarch):
    osarch = osarch.split(" ")
    filtered = []
    for arch in archs:
        if arch["GOOS"] + "/" + arch["GOARCH"] in osarch:
            filtered.append(arch)
    return filtered

if __name__ == '__main__':
    argparser = argparse.ArgumentParser()
    argparser.add_argument("-parallel", type=int, default=os.cpu_count())
    argparser.add_argument("-output", default="{{.Dir}}_{{.OS}}_{{.Arch}}")
    argparser.add_argument("-osarch", help="specify which pairs of os/architecture should be built. Example: -osarch='windows/arm64 linux/amd64")
    args, go_args = argparser.parse_known_args()
    archs = json.loads(subprocess.check_output("go tool dist list -json".split()))
    if args.osarch != None:
        archs = filter(archs, args.osarch)
    queue = Queue()
    for arch in archs:
        arch.update({"args": go_args, "template": args.output})
        queue.put(arch)
    print(f"Using {args.parallel} workers")
    for i in range(args.parallel):
        worker = Thread(target=build, args=(queue,))
        worker.setDaemon(True)
        worker.start()
    queue.join()