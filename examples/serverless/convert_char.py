# 将字符串中的每两位数字作为偏移量进行处理
# 若数字为 00，则转换为_
# 若数字介于 01-26，则转换为 a-z
# 若数字介于 27-52，则转换为 A-Z
# 若数字大于 53，则不做处理
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
                    offset = int(text[i:i+2])
                    if offset == 0:
                        result += "_"
                    elif 1 <= offset <= 26:
                        result += chr((offset - 1) % 26 + ord('a'))
                    elif 27 <= offset <= 52:
                        result += chr((offset - 27) % 26 + ord('A'))
                    else:
                        result += text[i:i+2]
                print(result)
        except ValueError:
            print(0)
