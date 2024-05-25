import sys

def add_one(number):
    return number + 1

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            number = int(sys.argv[1])
            result = add_one(number)
            print(result)
        except ValueError:
            print(0)