def read_in_chunks(file_object, chunk_size=1024):
    """Lazy function (generator) to read a file piece by piece.
    Default chunk size: 1k."""
    while True:
        data = file_object.read(chunk_size)
        if not data:
            break
        yield data

if __name__ == "__main__":
    path = "F:/GoWorkplace/goinaction/speechToText/audios/test.wav"
    with open(path, 'rb') as f:
        i = 0
        for piece in read_in_chunks(f, chunk_size=1024*16):
            with open(str(i) + ".txt", "wb") as tf:
                tf.write(piece)
            i += 1
            if i > 2:
                break
        