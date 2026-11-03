import json
import sys

from materialyoucolor.dynamiccolor.material_dynamic_colors import MaterialDynamicColors
from materialyoucolor.hct import Hct
from materialyoucolor.scheme.scheme_content import SchemeContent

accent, polarity = sys.argv[1], sys.argv[2]
scheme = SchemeContent(Hct.from_int(int(f"ff{accent}", 16)), polarity == "dark", 0.0)
colours = {c.name: f"{c.get_argb(scheme) & 0xFFFFFF:06x}" for c in MaterialDynamicColors().all_colors}

keys = [
    "primary", "onPrimary", "primaryContainer", "onPrimaryContainer", "inversePrimary",
    "primaryFixed", "primaryFixedDim", "onPrimaryFixed", "onPrimaryFixedVariant",
    "secondary", "onSecondary", "secondaryContainer", "onSecondaryContainer",
    "secondaryFixed", "secondaryFixedDim", "onSecondaryFixed", "onSecondaryFixedVariant",
]
out = {k: colours[k] for k in keys}
out["primary_paletteKeyColor"] = colours["primaryPaletteKeyColor"]
out["secondary_paletteKeyColor"] = colours["secondaryPaletteKeyColor"]
for k in "surfaceTint", "primary", "primaryFixedDim", "secondary", "secondaryFixedDim":
    out[k] = accent
print(json.dumps(out))
