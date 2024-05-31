import sys
import re

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            text = int(sys.argv[1])
            words = re.findall(r'\b\w+\b', text.lower())
            print(words)
        except ValueError:
            print(0)