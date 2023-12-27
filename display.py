from tkinter import Tk, Canvas, mainloop


x1=0
y1=1
x2=2
y2=3
r=4
g=5
b=6
if __name__ == "__main__":

    master = Tk()
    w = Canvas(master, width=800, height=800)
    w.pack()

    try:
        with open("color.out", 'r') as file:
            for line_number, line in enumerate(file, start=1):
                parts = line.strip().split(":")

                color = "#{}{}{}".format(
                    hex(int(parts[r]))[2:].zfill(2),
                    hex(int(parts[g]))[2:].zfill(2),
                    hex(int(parts[b]))[2:].zfill(2)
                )

                w.create_rectangle(
                    int(parts[x1]), 
                    int(parts[y1]), 
                    int(parts[x2]), 
                    int(parts[y2]), 
                    fill=color, 
                    outline=color
                )
    except FileNotFoundError:
        print(f"Error: File '{file_path}' not found.")
    except Exception as e:
        print(f"Error: {e}")


    mainloop()