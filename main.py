import tkinter as tk
from tkinter import ttk, messagebox
import json
from pathlib import Path


class Library:
    def __init__(self, root):
        self.root = root
        self.root.title("my abyss")
        self.geometry("900x700")
        self.configure("#BDB76B")

        # data configures
        self.books = []

        # UI configures





def main():
    root = tk.Tk()
    lib = Library(root)
    root.mainloop()
