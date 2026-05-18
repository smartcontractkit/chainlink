COVENANT_VERSION = "IUSC-1.0"
COVENANT_HASH = "8f5e...9b2a" # The hash of your Ironclad Covenant
IDENTITY = "The Keeper / Jon S. (Pray4Love1)"

def get_sovereign_prefix():
    return f"[{COVENANT_VERSION}:{COVENANT_HASH}] AUTHORIZED BY {IDENTITY}"
