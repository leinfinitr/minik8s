import sys

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            file_name = int(sys.argv[1])
            with open(file_name, 'r') as f:
                text = f.read()
            print(text)
        except ValueError:
            print(0)