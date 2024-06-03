import sys
import os

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            file_path = os.path.abspath(sys.argv[1])
            with open(file_path, "r") as f:
                string = f.read()
            result = ""
            for char in string:
                if char.isdigit():
                    result += str((int(char) + 7) % 10)
                else:
                    result += char
            print(result)
        except ValueError:
            print(1)