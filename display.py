import argparse

parser = argparse.ArgumentParser(description="Read custom pixel coordinate file")
parser.add_argument("-f","--file",type=str,help="Input path to coordinate file")
args = parser.parse_args()

if __name__ == "__main__":

    try:
        file_path = args.file

        if file_path is None:
            raise Exception("No file specified")
        
        from tkinter import Tk, Canvas, mainloop

        X1=HEIGHT=0
        Y1=WIDTH=1
        X2=2
        Y2=3
        R=4
        G=5
        B=6

        with open(file_path, 'r') as file:
            h_w = file.readline().strip().split(":")

            master = Tk()
            master.title("Decoded JPEG")
            w = Canvas(master, width=int(h_w[WIDTH]), height=int(h_w[HEIGHT]))
            w.pack()

            for line_number, line in enumerate(file, start=1):
                parts = line.strip().split(":")

                color = "#{}{}{}".format(
                    hex(int(parts[R]))[2:].zfill(2),
                    hex(int(parts[G]))[2:].zfill(2),
                    hex(int(parts[B]))[2:].zfill(2)
                )

                w.create_rectangle(
                    int(parts[X1]), 
                    int(parts[Y1]), 
                    int(parts[X2]), 
                    int(parts[Y2]), 
                    fill=color,
                    outline=color
                )
           
            mainloop()
    except FileNotFoundError:
        print(f"Error: File '{file_path}' not found.")
    except Exception as e:
        print(f"Error: {e}")