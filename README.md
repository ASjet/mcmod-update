# mcmod-update

A commandline tool for updating your Minecraft mods from [CurseForge](https://www.curseforge.com/minecraft).

## Usage

First, find the *Project ID* and *name* of the mod that you would like to download from [CurseForge Mod](https://www.curseforge.com/minecraft/search?page=1&class=mc-mods&sortType=1&pageSize=20), and put them in a file, separate with comma, one mod a line, like this:

```txt
238222,JET
248787,Apple Skin
```

Move the file, say `mods.txt`, to `.minecraft/mods`, and run the command:

```bash
# All following command should be run with the environment variable CurseForgeAPIKey set
export CurseForgeAPIKey=${your_curseforge_api_key}

cat mods.txt | modver > ver.json
```

Then you will get a `ver.json` file, which contains the meta info of the mods you want to download.

Finally, run the command:

```bash
modver -c ver.json -d
```

and the mods will be downloaded to current directory.

## Update mods

Navigate to the `.minecraft/mods` directory, and run the command:

```bash
modver -c ver.json
```

And you will see if there is any update. If it is, run the command:

```bash
modver -c ver.json -d
```

that will download the latest mods and replace the old ones.
