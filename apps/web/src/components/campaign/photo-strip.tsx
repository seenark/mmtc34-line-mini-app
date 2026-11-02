import { cn } from "@mmtc34/ui/lib/utils";
import { useRef, useState } from "react";

import { imageUrl } from "@/lib/image";

export function PhotoStrip({ photos }: { photos: string[] }) {
  const stripRef = useRef<HTMLDivElement>(null);
  const [active, setActive] = useState(0);

  if (photos.length === 0) return null;

  const onScroll = () => {
    const strip = stripRef.current;
    if (!strip) return;
    setActive(Math.round(strip.scrollLeft / strip.clientWidth));
  };

  const scrollTo = (index: number) => {
    const strip = stripRef.current;
    if (!strip) return;
    const reduceMotion = window.matchMedia(
      "(prefers-reduced-motion: reduce)",
    ).matches;
    strip.scrollTo({
      left: index * strip.clientWidth,
      behavior: reduceMotion ? "auto" : "smooth",
    });
  };

  return (
    <section
      aria-roledescription="carousel"
      aria-label="รูปภาพอีเวนต์"
      className="flex flex-col gap-2"
    >
      <div
        ref={stripRef}
        onScroll={onScroll}
        className="flex snap-x snap-mandatory [scrollbar-width:none] overflow-x-auto overscroll-x-contain rounded-xl [&::-webkit-scrollbar]:hidden"
      >
        {photos.map((photo, i) => (
          <div key={photo} className="w-full shrink-0 snap-center">
            <img
              src={imageUrl(photo)}
              alt=""
              loading={i === 0 ? "eager" : "lazy"}
              className="bg-muted aspect-square max-h-80 w-full rounded-xl object-contain"
            />
          </div>
        ))}
      </div>
      {photos.length >= 2 && (
        <div className="flex justify-center">
          {photos.map((photo, i) => (
            <button
              key={photo}
              type="button"
              aria-label={`รูปที่ ${i + 1}`}
              aria-current={i === active}
              onClick={() => scrollTo(i)}
              className="focus-visible:ring-ring rounded-full p-2 outline-none focus-visible:ring-2"
            >
              <span
                className={cn(
                  "block size-2 rounded-full transition-colors",
                  i === active ? "bg-primary" : "bg-muted-foreground/30",
                )}
              />
            </button>
          ))}
        </div>
      )}
    </section>
  );
}
