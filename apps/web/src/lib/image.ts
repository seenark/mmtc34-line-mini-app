const MAX_SIDE = 1600;
const QUALITY = 0.82;

/** Downscales to MAX_SIDE and re-encodes as JPEG so phone photos stay under the server's 5 MB cap. */
export async function shrinkImage(file: File): Promise<Blob> {
  const bitmap = await createImageBitmap(file, {
    imageOrientation: "from-image",
  });
  try {
    const scale = Math.min(1, MAX_SIDE / Math.max(bitmap.width, bitmap.height));
    const canvas = document.createElement("canvas");
    canvas.width = Math.round(bitmap.width * scale);
    canvas.height = Math.round(bitmap.height * scale);
    const ctx = canvas.getContext("2d");
    if (!ctx) throw new Error("ไม่สามารถประมวลผลรูปภาพได้");
    ctx.drawImage(bitmap, 0, 0, canvas.width, canvas.height);
    return await new Promise<Blob>((resolve, reject) =>
      canvas.toBlob(
        (blob) =>
          blob ? resolve(blob) : reject(new Error("ไม่สามารถประมวลผลรูปภาพได้")),
        "image/jpeg",
        QUALITY,
      ),
    );
  } finally {
    bitmap.close();
  }
}

/** Absolute-path URL of an uploaded image; name is "" when none is set. */
export const imageUrl = (name: string) =>
  `${import.meta.env.BASE_URL}api/images/${name}`;
