import sys

def add_seven(number):
    return number - 7

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            number = int(sys.argv[1])
            result = add_seven(number)
            print(result)
        except ValueError:
            print(0)