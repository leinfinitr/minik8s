# 将字符串中的每两位数字作为一个字母的ASCII码，转换为字母
# 01-26 -> a-z
# 27-52 -> A-Z
import sys

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            text = sys.argv[1]
            if len(text) % 2 != 0:
                print(0)
            else:
                result = ""
                for i in range(0, len(text), 2):
                    num = int(text[i:i+2])
                    if num < 1 or num > 52:
                        print(0)
                        break
                    if num <= 26:
                        result += chr(num + 96)
                    else:
                        result += chr(num + 38)
                print(result)
        except ValueError:
            print(0)
