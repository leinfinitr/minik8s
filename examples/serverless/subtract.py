# 将一个字符串中的每个数字减7的个位数替换该数字，然后打印替换后的字符串
import sys

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            number = sys.argv[1]
            result = ""
            for char in number:
                result += str((int(char) - 7) % 10)
            print(result)
        except ValueError:
            print(0)
