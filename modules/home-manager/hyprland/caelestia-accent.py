import json
import sys

from materialyoucolor.dynamiccolor.material_dynamic_colors import MaterialDynamicColors
from materialyoucolor.hct import Hct
from materialyoucolor.scheme.scheme_content import SchemeContent

accent, tertiary, polarity = sys.argv[1], sys.argv[2], sys.argv[3]

keys = [
    "primary", "onPrimary", "primaryContainer", "onPrimaryContainer", "inversePrimary",
    "primaryFixed", "primaryFixedDim", "onPrimaryFixed", "onPrimaryFixedVariant",
    "secondary", "onSecondary", "secondaryContainer", "onSecondaryContainer",
    "secondaryFixed", "secondaryFixedDim", "onSecondaryFixed", "onSecondaryFixedVariant",
]


def family(hex):
    scheme = SchemeContent(Hct.from_int(int(f"ff{hex}", 16)), polarity == "dark", 0.0)
    colours = {c.name: f"{c.get_argb(scheme) & 0xFFFFFF:06x}" for c in MaterialDynamicColors().all_colors}
    out = {k: colours[k] for k in keys}
    out["primary_paletteKeyColor"] = colours["primaryPaletteKeyColor"]
    out["secondary_paletteKeyColor"] = colours["secondaryPaletteKeyColor"]
    for k in "primary", "primaryFixedDim", "secondary", "secondaryFixedDim":
        out[k] = hex
    return out


out = family(accent)
out["surfaceTint"] = accent
for k, v in family(tertiary).items():
    if "rimary" in k:
        out[k.replace("primary", "tertiary").replace("Primary", "Tertiary")] = v
print(json.dumps(out))
