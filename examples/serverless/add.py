# 将一个字符串中的每个数字加7的个位数替换该数字，然后打印替换后的字符串
import sys

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            number = sys.argv[1]
            result = ""
            for char in number:
                if char.isdigit():
                    result += str((int(char) + 7) % 10)
                else:
                    result += char
            print(result)
        except ValueError:
            print(0)
