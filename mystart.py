#!/usr/bin/env python3

import subprocess
from math import trunc
import mystart_data
from colorama import Fore
from prettytable import PrettyTable


class MyStartScraper:
    def __init__(self):
        self.vars = {}
        self.vars["messages"] = []

    def vars_generate(self):
        self.vars["uname"] = subprocess.getoutput(
            ["cat /etc/os-release | grep \"PRETTY\" | cut -d'=' -f2-"]
        )
        self.vars["distro"] = subprocess.getoutput(
            ['grep -Po "(?<=^ID=).+" /etc/os-release | sed \'s/"//g\'']
        )
        self.vars["host"] = subprocess.getoutput(["uname -n"])
        self.vars["cpu_cores"] = int(
            subprocess.getoutput(
                [
                    "cat /proc/cpuinfo | grep \"cpu cores\" | head -n 1 | awk '{print $4}'"
                ]
            )
        )
        checkip4 = subprocess.getoutput(["ip route get 8.8.8.8 2>/dev/null"])
        if not checkip4:
            self.vars["ipv4"] = "N/A"
        else:
            self.vars["ipv4"] = subprocess.getoutput(
                ["ip route get 8.8.8.8 | grep src | awk '{print $7}'"]
            )

        checkip6 = subprocess.getoutput(
            ["ip route get 2001:4860:4860::8888 2>/dev/null"]
        )
        if not checkip6:
            self.vars["ipv6"] = "N/A"
        else:
            self.vars["ipv6"] = subprocess.getoutput(
                ["ip route get 2001:4860:4860::8888 | grep src | awk '{print $11}'"]
            )

        self.vars["lastlog"] = subprocess.getoutput(
            [
                'last | head -n 2 | tail -1 | awk \'{print $1 " on " $2 ", " $4 " " $5 " " $6 " " $7 " from " $3}\''
            ]
        )
        self.vars["uptime"] = subprocess.getoutput(["cut -d. -f1 /proc/uptime"])
        self.vars["uptime"] = int(self.vars["uptime"])
        self.vars["up_days"] = trunc((self.vars.get("uptime") / 60 / 60 / 24))
        self.vars["up_hours"] = trunc((self.vars.get("uptime") / 60 / 60 % 24))
        self.vars["up_minutes"] = trunc((self.vars.get("uptime") / 60 % 60))
        self.vars["up_seconds"] = self.vars.get("uptime") % 60
        self.vars["loadavg"] = subprocess.getoutput(
            ['cat /proc/loadavg | awk \'{print $1 " " $2 " " $3 " " $4}\'']
        )

        cpu_speed_cmd1 = subprocess.run(
            [
                "cat",
                "/proc/cpuinfo",
            ],
            check=True,
            text=True,
            stdout=subprocess.PIPE,
        )
        cpu_speed_cmd1 = cpu_speed_cmd1.stdout.strip("\n")
        cpu_speed_cmd2 = subprocess.run(
            [
                "grep",
                "MHz",
            ],
            input=cpu_speed_cmd1,
            check=True,
            text=True,
            stdout=subprocess.PIPE,
        )
        cpu_speed_cmd2 = cpu_speed_cmd2.stdout.strip("\n")
        cpu_speed_cmd3 = subprocess.run(
            ["awk", "{print $4}"],
            input=cpu_speed_cmd2,
            capture_output=True,
            check=True,
            text=True,
        )
        cpu_speeds = 0.0
        thread_count = 0
        for line in cpu_speed_cmd3.stdout.splitlines():
            thread_count += 1
            cpu_speeds += float(line)

        self.vars["cpu_hz"] = round(((cpu_speeds) / (thread_count)) / 1000, 2)
        self.vars["cpu_threads"] = thread_count
        cpu_used = float(
            subprocess.getoutput(["ps -eo pcpu | awk '{tot=tot+$1} END {print tot}'"])
        )

        cpu_usage = round((cpu_used) / (self.vars.get("cpu_threads")), 2)
        self.vars["cpu_usage"] = f"{cpu_usage} %"
        memory_total = float(
            subprocess.getoutput(
                ["cat /proc/meminfo | grep MemTotal | awk '{print $2}'"]
            )
        )
        memory_avail = float(
            subprocess.getoutput(
                ["cat /proc/meminfo | grep MemAvailable | awk '{print $2}'"]
            )
        )
        memtot = round(((memory_total / 1024) / 1024), 2)
        self.vars["memtot"] = f"{memtot}G"
        memory_used = round(((memory_total - memory_avail) / 1024 / 1024), 2)
        self.vars["memuse"] = f"{memory_used}G"

        swap_total = float(
            subprocess.getoutput(
                ["cat /proc/meminfo | grep SwapTotal | awk '{print $2}'"]
            )
        )
        swap_free = float(
            subprocess.getoutput(
                ["cat /proc/meminfo | grep SwapFree | awk '{print $2}'"]
            )
        )
        swaptot = round(((swap_total / 1024) / 1024), 2)
        self.vars["swaptot"] = f"{swaptot}G"
        swap_used = round(((swap_total - swap_free) / 1024 / 1024), 2)
        self.vars["swapuse"] = f"{swap_used}G"

        self.vars["diskuse"] = subprocess.getoutput(
            ["df -h | awk '{if($(NF) == \"/\") {print $(NF-1); exit;}}'"]
        )
        self.vars["disksize"] = subprocess.getoutput(
            ["df -h | awk '{if($(NF) == \"/\") {print $(NF-4); exit;}}'"]
        )
        try:
            disk_cmd1 = subprocess.run(
                [
                    "df",
                    "-h",
                ],
                check=True,
                text=True,
                stdout=subprocess.PIPE,
            ).stdout.strip("\n")

            disk_cmd2 = subprocess.run(
                [
                    "grep",
                    "sd",
                ],
                input=disk_cmd1,
                check=True,
                text=True,
                stdout=subprocess.PIPE,
            ).stdout.strip("\n")

            disk_total = 0.0
            disk_used = 0.0
            for line in disk_cmd2.splitlines():
                disk_total += float(
                    subprocess.run(
                        [
                            "awk",
                            "{print $2}",
                        ],
                        input=line,
                        check=True,
                        text=True,
                        capture_output=True,
                    ).stdout.strip("\n")[:-1]
                )

                disk_used += float(
                    subprocess.run(
                        [
                            "awk",
                            "{print $3}",
                        ],
                        input=line,
                        check=True,
                        text=True,
                        capture_output=True,
                    ).stdout.strip("\n")[:-1]
                )

            if not disk_total:
                self.vars["disk_pool_size"] = "N/A"
            else:
                self.vars["disk_pool_size"] = f"{round(disk_total, 1)}TB"

            if not disk_used:
                self.vars["disk_pool_used"] = "N/A"
            else:
                self.vars["disk_pool_used"] = f"{round(disk_used, 1)}TB"
        except:
            self.vars["disk_pool_size"] = "N/A"
            self.vars["disk_pool_used"] = "N/A"

        try:
            test = subprocess.run(
                ["sudo", "liquidctl", "status"],
                capture_output=True,
                check=True,
                text=True,
                timeout=2,
            )
            self.vars["front_left_fan"] = subprocess.getoutput(
                ["sudo liquidctl status | grep \"Fan 1\" | awk '{print $5}'"]
            )
            self.vars["front_right_fan"] = subprocess.getoutput(
                ["sudo liquidctl status | grep \"Fan 2\" | awk '{print $5}'"]
            )
            self.vars["back_fan"] = subprocess.getoutput(
                ["sudo liquidctl status | grep \"Fan 3\" | awk '{print $5}'"]
            )
            self.vars["cpu_temp"] = subprocess.getoutput(
                ["sudo liquidctl status | grep \"Temp sensor 1\" | awk '{print $5}'"]
            )
            self.vars["memory_temp"] = subprocess.getoutput(
                ["sudo liquidctl status | grep \"Temp sensor 2\" | awk '{print $5}'"]
            )
            self.vars["drive_1_temp"] = subprocess.getoutput(
                ["sudo liquidctl status | grep \"Temp sensor 3\" | awk '{print $5}'"]
            )
            self.vars["drive_2_temp"] = subprocess.getoutput(
                ["sudo liquidctl status | grep \"Temp sensor 4\" | awk '{print $5}'"]
            )
        except:
            self.vars["front_left_fan"] = "N/A"
            self.vars["front_right_fan"] = "N/A"
            self.vars["back_fan"] = "N/A"
            self.vars["cpu_temp"] = "N/A"
            self.vars["memory_temp"] = "N/A"
            self.vars["drive_1_temp"] = "N/A"
            self.vars["drive_2_temp"] = "N/A"

        trans_kick_status = subprocess.getoutput(
            [
                "systemctl status trans-kick.service | grep Active | awk '{ print $2, $3 }'"
            ]
        )

        if trans_kick_status == "active (running)":
            self.vars[
                "trans_kick_status"
            ] = f"{trans_kick_status} {mystart_data.thumb_up}"
        elif trans_kick_status == "Unit trans-kick.service could not be found.":
            self.vars["trans_kick_status"] = "N/A"
        else:
            self.vars[
                "trans_kick_status"
            ] = f"{trans_kick_status} {mystart_data.stop_emoji}"

        try:
            self.vars["nord_addr"] = subprocess.run(
                ["docker", "exec", "nord", "curl", "ifconfig.io"],
                capture_output=True,
                check=True,
                text=True,
                timeout=2,
            ).stdout.strip("\n")
        except:
            self.vars["nord_addr"] = "N/A"
        try:
            self.vars["trans_addr"] = subprocess.run(
                ["docker", "exec", "transmission", "curl", "ifconfig.io"],
                capture_output=True,
                check=True,
                text=True,
                timeout=2,
            ).stdout.strip("\n")
        except:
            self.vars["trans_addr"] = "N/A"

        if self.vars["trans_kick_status"] == "N/A":
            self.vars["vpn_check"] = "N/A"

        elif (
            self.vars["nord_addr"] == self.vars["trans_addr"]
            and self.vars["nord_addr"] != "N/A"
        ):
            self.vars["vpn_check"] = f"Protected {mystart_data.thumb_up}"
        else:
            self.vars["vpn_check"] = f"Unprotected {mystart_data.stop_emoji}"
        self.vars["thislog"] = subprocess.getoutput(
            [
                'lastlog -u $USER | tail -n 1 | awk \'{print $4 " " $5 " " $6 " " $7 " from " $3}\''
            ]
        )
        self.vars["user"] = subprocess.getoutput(["id -un"])
        self.vars["psu"] = subprocess.getoutput(["ps -aux | grep -i $USER | wc -l"])
        self.vars["psa"] = subprocess.getoutput(["ps -aux | wc -l"])
        self.vars["users"] = subprocess.getoutput(['who -q | tail -n1 | cut -c "9-11"'])
        self.vars["sshusers"] = subprocess.getoutput(["ss | grep -i ssh | wc -l"])

        line_border = f"{Fore.MAGENTA}======================================================================={Fore.RESET}"
        self.vars["messages"].insert(0, line_border)
        self.vars["messages"].insert(2, line_border)
        msg2 = (
            f"{Fore.GREEN}[*]{Fore.RESET} System details\t\t:{Fore.MAGENTA} %s | %s |"
        ) % (self.vars.get("uname"), self.vars.get("distro"))
        self.vars["messages"].insert(4, msg2)
        msg3 = (
            f"{Fore.GREEN}[*]{Fore.RESET} System uptime\t\t:{Fore.MAGENTA} %s days %s hours %s minutes %s seconds"
        ) % (
            self.vars.get("up_days"),
            self.vars.get("up_hours"),
            self.vars.get("up_minutes"),
            self.vars.get("up_seconds"),
        )
        self.vars["messages"].insert(5, msg3)
        msg4 = (f"{Fore.GREEN}[*]{Fore.RESET} System load\t\t:{Fore.MAGENTA} %s") % (
            self.vars.get("loadavg")
        )
        self.vars["messages"].insert(6, msg4)
        msg5 = (
            f"{Fore.GREEN}[*]{Fore.RESET} CPU info\t\t\t:{Fore.MAGENTA} %s in use of %scores/%sthreads at %sGHz"
        ) % (
            self.vars.get("cpu_usage"),
            self.vars.get("cpu_cores"),
            self.vars.get("cpu_threads"),
            self.vars.get("cpu_hz"),
        )
        self.vars["messages"].insert(7, msg5)
        msg6 = (
            f"{Fore.GREEN}[*]{Fore.RESET} Memory in use\t\t:{Fore.MAGENTA} %s of %s"
        ) % (self.vars.get("memuse"), self.vars.get("memtot"))
        self.vars["messages"].insert(8, msg6)
        msg7 = (
            f"{Fore.GREEN}[*]{Fore.RESET} Swap memory in use\t\t:{Fore.MAGENTA} %s of %s"
        ) % (self.vars.get("swapuse"), self.vars.get("swaptot"))
        self.vars["messages"].insert(9, msg7)
        msg8 = (
            f"{Fore.GREEN}[*]{Fore.RESET} Root disk usage\t\t:{Fore.MAGENTA} %s of %s"
        ) % (self.vars.get("diskuse"), self.vars.get("disksize"))
        self.vars["messages"].insert(10, msg8)
        msg9 = (
            f"{Fore.GREEN}[*]{Fore.RESET} Disk pool available\t:{Fore.MAGENTA} %s"
        ) % (self.vars.get("disk_pool_size"))
        self.vars["messages"].insert(11, msg9)
        msg10 = (
            f"{Fore.GREEN}[*]{Fore.RESET} Disk pool used\t\t:{Fore.MAGENTA} %s"
        ) % (self.vars.get("disk_pool_used"))
        self.vars["messages"].insert(12, msg10)
        msg14 = (
            f"{Fore.GREEN}[*]{Fore.RESET} Last system login\t\t:{Fore.MAGENTA} %s"
        ) % (self.vars.get("lastlog"))
        self.vars["messages"].insert(16, msg14)

        msg20 = (f"{Fore.GREEN}[*]{Fore.RESET} IPv4 address\t\t:{Fore.MAGENTA} %s") % (
            self.vars.get("ipv4")
        )
        self.vars["messages"].insert(21, line_border)
        self.vars["messages"].insert(22, msg20)
        msg21 = (f"{Fore.GREEN}[*]{Fore.RESET} IPv6 address\t\t:{Fore.MAGENTA} %s") % (
            self.vars.get("ipv6")
        )
        self.vars["messages"].insert(23, msg21)
        self.vars["messages"].insert(24, line_border)
        welcome = (
            f"{Fore.WHITE}Hello %s and welcome to %s {mystart_data.grinning_face}"
        ) % (self.vars.get("user"), self.vars.get("host"))

        self.vars["messages"].insert(1, welcome)
        msg1 = (f"{Fore.GREEN}[*]{Fore.RESET} Login details\t\t:{Fore.MAGENTA} %s") % (
            self.vars.get("thislog")
        )
        self.vars["messages"].insert(3, msg1)
        msg11 = (
            f"{Fore.GREEN}[*]{Fore.RESET} System processes\t\t:{Fore.MAGENTA} %s running %s, total of %s running on %s"
        ) % (
            self.vars.get("user"),
            self.vars.get("psu"),
            self.vars.get("psa"),
            self.vars.get("host"),
        )
        self.vars["messages"].insert(13, msg11)
        msg12 = (
            f"{Fore.GREEN}[*]{Fore.RESET} SSH logins\t\t\t:{Fore.MAGENTA} %s user(s) logged in via SSH"
        ) % (self.vars.get("sshusers"))
        self.vars["messages"].insert(14, msg12)
        msg13 = (
            f"{Fore.GREEN}[*]{Fore.RESET} Total logins\t\t:{Fore.MAGENTA} %s user(s) in total logged in"
        ) % (self.vars.get("users"))
        self.vars["messages"].insert(15, msg13)

        if self.vars["host"] == "titan":
            msg16 = (
                f"{Fore.GREEN}[*]{Fore.RESET} Left drive bay fan\t\t:{Fore.MAGENTA} %s/RPM with bay temp: %s/C"
            ) % (self.vars.get("front_left_fan"), self.vars.get("drive_1_temp"))
            self.vars["messages"].insert(17, line_border)
            self.vars["messages"].insert(18, msg16)
            msg17 = (
                f"{Fore.GREEN}[*]{Fore.RESET} Right drive bay fan\t:{Fore.MAGENTA} %s/RPM with bay temp: %s/C"
            ) % (self.vars.get("front_right_fan"), self.vars.get("drive_2_temp"))
            self.vars["messages"].insert(19, msg17)
            msg18 = (
                f"{Fore.GREEN}[*]{Fore.RESET} Back chassis fan\t\t:{Fore.MAGENTA} %s/RPM with Memory temp: %s/C & CPU temp: %s/C"
            ) % (
                self.vars.get("back_fan"),
                self.vars.get("memory_temp"),
                self.vars.get("cpu_temp"),
            )
            self.vars["messages"].insert(20, msg18)

        if self.vars["user"] == "orion" and self.vars["host"] == "titan":
            msg22 = (
                f"{Fore.GREEN}[*]{Fore.RESET} Nord address\t\t:{Fore.MAGENTA} %s"
            ) % (self.vars.get("nord_addr"))
            self.vars["messages"].insert(25, msg22)
            msg23 = (
                f"{Fore.GREEN}[*]{Fore.RESET} Transmission address\t:{Fore.MAGENTA} %s"
            ) % (self.vars.get("trans_addr"))
            self.vars["messages"].insert(26, msg23)
            self.vars["messages"].insert(27, line_border)
            msg24 = (
                f"{Fore.GREEN}[*]{Fore.RESET} Transmission status\t:{Fore.MAGENTA} %s"
            ) % (self.vars.get("vpn_check"))
            self.vars["messages"].insert(28, msg24)
            msg25 = (
                f"{Fore.GREEN}[*]{Fore.RESET} Trans-kick status\t\t:{Fore.MAGENTA} %s"
            ) % (self.vars.get("trans_kick_status"))
            self.vars["messages"].insert(29, msg25)
            self.vars["messages"].insert(30, line_border)
        return

    def payload_table(self):
        table = PrettyTable()
        table.border = False
        table.field_names = [self.vars["messages"][0]]

        for i in self.vars["messages"][1:]:
            table.add_row([i])
        table.align = "l"

        print(table)
        return


if __name__ == "__main__":
    scraper = MyStartScraper()
    scraper.vars_generate()
    scraper.payload_table()
