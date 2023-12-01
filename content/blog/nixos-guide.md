+++
title = "How to get started with NixOS"
date = 2022-11-23

[taxonomies]
series=["NixOS Desktop"]
tags=["Nix", "Linux"]
+++

# How to get started with NixOS [WIP]

Currently, NixOS isn't as well documented as arch and many guides that are
already out there are either outdated or incorrect. The fact that you have
to learn an entire new functional language just to be able to properly use
your system, also doesn't help with the learning curve.

That's why I made this guide.

<!-- more -->

# Why would you use NixOS?

First of all, you need to ask yourself how would you benefit from using NixOS,
or why would you switch

[Too long; Didn't read](#tldr)

## _Declarative, reproducible and immutable_ {#propaganda}

The very first thing you have probably heard about NixOS is that its _declarative,
reproducible and immutable_. But what does it even mean?

Basically _everything_ about your system is declared in your configuration, from the [kernel](https://git.sioodmy.dev/sioodmy/dotfiles/src/commit/48edd68feafa442a1a9bf0ce133d71419977dbe7/modules/core/bootloader.nix#L36)
to even [firefox addons](https://git.sioodmy.dev/sioodmy/dotfiles/src/commit/48edd68feafa442a1a9bf0ce133d71419977dbe7/modules/home/schizofox/default.nix#L158)

Its literally a dotfiles repo, but on anabolic steroids.

Sounds cool? Wait till you have to install anything from source.

> hint: running `sudo make install` won't work, it would be too easy

You need to write a _derivation_ or at least an _overlay_.

Okay, but why should you care?

Let's say you just bought a new drive and want to reinstall your system on it.
All you have to do is just to clone your repo, run `nixos-install` and wait.

Not cool enough?

Your school forces you to use windows on their computers, but you're used to neovim, your
shell aliases, etc? You can make a portable [NixOS WSL](https://github.com/nix-community/NixOS-WSL) tarball

Still not cool enough?

Imagine making a custom ISO with your configuration, you plug it in anywhere
and because its reproducible it would just work

# pkgs

Nix is the most advanced package manager on any existing system.

Want to use _libressl_ for your web server, but don't want to break any packages
using openssl as a dependency? No problem on nix.

```nix
  nixpkgs.overlays = [
    (final: super: {
      nginxStable = super.nginxStable.override {openssl = super.pkgs.libressl;};
    })
  ];
```

According to [this graph](https://repology.org/graph/map_repo_size_fresh.svg) nixpkgs is currently the biggest linux repo out there.

<img src="https://repology.org/graph/map_repo_size_fresh.svg" width="100%" alt="graph" />

> Actually those stats are quite boosted, because it counts every single derivation
> as a package

Ok, but let's say you want to install _picom-ibhagwan-git_ as you would usually do
from the AUR, but wait its not there. What should I do?

That's when we can simply override existing _picom_ package, because the installation steps
are quite simliar

```nix
  package = pkgs.picom.overrideAttrs(o: {
        src = pkgs.fetchFromGitHub {
          repo = "picom";
          owner = "ibhagwan";
          rev = "44b4970f70d6b23759a61a2b94d9bfb4351b41b1";
           sha256 = "0iff4bwpc00xbjad0m000midslgx12aihs33mdvfckr75r114ylh";
        };
  });
```

### tl;dr {#tldr}

**Nix advantages**:

- Correct and complete packaging
- Immutable & reproducible results
- Easy to cross and static compile
- Source-based (you can alter packages without forking anything)
- Single package manager to rule them all! (C, Python, Docker, NodeJS, etc)
- Great for development, easily switches between dev envs with direnv
- Easy to set up a binary cache
- By default has a binary cache so you almost never need to compile anything
- Easy to set up remote building
- Excellent testing infrastructure
- Portable - runs on Linux and macOS
- Can be built statically and run anywhere without root permissions
- Mix and match different package versions without conflicts

**NixOS advantages**:

- Declarative configuration
- Easy to deploy machines and their configuration
- Out of the box Rollbacks.
- Configuration options for many programs & services
- Free of side effects - Actually uninstalls packages and their dependencies
- Easy to set up VMs

Now that you know why you would want to use Nix and NixOS, we can finally begin with installation

# Bootable usb

```bash
$ curl -fL https://channels.nixos.org/nixos-22.05/latest-nixos-gnome-x86_64-linux.iso -o nixos.iso
# make a bootable usb (/dev/sdX)
$ sudo cp nixos.iso /dev/sdX
$ sync
```

# First steps

Just boot from your usb and you should end up with screen looking like this.

<img src="https://media.discordapp.net/attachments/979483369958150174/1045429395852103690/image.png?width=1118&height=629" alt="installer" width="100%" />

Close the calamares installer and open a terminal like a true chad.

As you may already guessed, now you need to partition your hard drive.
Just make a _boot, swap and root_ partition with a tool of your choice.

I'm gonna use cfdisk, but it doesnt matter.

```bash
$ lsblk
$ sudo cfdisk /dev/sdX
# I prefer to use labels, rather than UUIDs
$ sudo mkfs.vfat -n boot
$ sudo mkswap /dev/sdX2 -L swap
$ sudo mkfs.ext4 /dev/sdX3 -L root
```

Now mount everything, you know how it goes.

```bash
$ sudo mount /dev/disk/by-label/root /mnt
$ sudo mkdir /mnt/boot
$ sudo mount /dev/disk/by-label/boot /mnt/boot
$ sudo swapon /dev/disk/by-label/swap
```

# Basic configuration

This is when the fun part begins. Start by generating default NixOS config like this

```bash
sudo nixos-generate-config --root /mnt
cd /mnt
```

Command above should generate _configuration.nix_ and _hardware-configuration_ in your /mnt/etc/nixos directory.
In theory, at this point you can run `nixos-install`, but we are going to tweak it a little.

```bash
sudo vim /mnt/etc/nixos/configuration.nix
```

The setup has about ~100 LoC and is well commented, so I highly recommend reading it all and making any desired changes.

Here are some things that you might want to change. Remove everything thats not necessary to clean it up.

### User configuration

```nix
users.user.yourname = {
    isNomalUser = true;
    groups = ["wheel"];
    initialPassword = "ChangeMe";
    packages = with pkgs; [
        firefox
        neovim
        ...
    ];
};
```

### System packages

Hint: Use [Nix search](https://search.nixos.org/packages#) to lookup package names

```nix
environment.systemPackages = with pkgs; [
  git # make sure you have git installed
];
```

### Systemd-boot

```nix
boot.loader = {
    systemd-boot.enable = true;
    efi.canTouchEfiVariables = true;
    timeout = 2; # you can also customize bootloader timeout here
};
```

### Sound

Not sure why thats not enabled by default btw.

Enable pulseaudio...

```nix
sound.enable = true;
hardware.pulseaudio.enable = true;
```

... or even better, pipewire

```nix
sound = {
   enable = true;
   mediaKeys.enable = true;
};

hardware.pulseaudio.enable = false;

services.pipewire = {
     enable = true;
     alsa = {
       enable = true;
       support32Bit = true;
     };
     wireplumber.enable = true;
     pulse.enable = true;
     jack.enable = true;
};
```

# Finally

Keep your fingers crossed and run

```nix
sudo nixos-install
```

<img src="https://media.discordapp.net/attachments/984919316250132540/1045757418849701969/image.png?width=1118&height=629" alt="Installation screenshot" width="100%"/>

Now we wait, take a break, drink water, go outside. It took about 20 minutes on my vm.

After its done set a root password and you can finally reboot.

- tip: If your internet connection gets interrupted try running following command:

```bash
$ nixos-rebuild switch --option substitute false
```

# Flakes

I couldn't find any easy to understand definition of _flake_, but all you really need to know (at least for now)
is that flake is that dotfiles repo, we talked about [earlier](#propaganda)

We will use [colemickens'](https://github.com/colemickens) nixos flake example, just to get started.

```bash
# you have git installed, right?
$ git clone https://github.com/colemickens/nixos-flake-example dotfiles
$ cd dotfiles
$ rm ./configuration.nix # we are going to replace it
$ nvim flake.nix
```

As you can see it's pretty minimal (note that I removed lines with "ignore" comment)

```nix
{
  description = "An example NixOS configuration";

  inputs = {
    nixpkgs = { url = "github:nixos/nixpkgs/nixos-unstable"; };
    nur = { url = "github:nix-community/NUR"; };
  };

  outputs = inputs:
  {
    nixosConfigurations = {
      mysystem = inputs.nixpkgs.lib.nixosSystem {
        system = "x86_64-linux";
        modules = [
          ./configuration.nix
        ];
        specialArgs = { inherit inputs; };
      };
    };
  };
}
```

Replace `mysystem` with your hostname

and we can also remove the NUR input

```diff
- nur = { url = "github:nix-community/NUR"; };
```

Remember the configuration.nix file we modified?
Copy it to the root directory of our flake along with the hardware config

```bash
$ cp -v /etc/nixos/configuration.nix /etc/nixos/hardware-configuration.nix .
$ git add . # this is required for nix to see the hardware configuration we just copied
```

The problem is that configuration.nix calls hardware-configuration.nix directly, which is quite dumb, since we want different hardware config
for each host. That's why we're removing the import

```diff
-  imports =
-    [ # Include the results of the hardware scan.
-      ./hardware-configuration.nix
-    ];
```

And add it as a module

```diff
yourhostname = inputs.nixpkgs.lib.nixosSystem {
  system = "x86_64-linux";
  modules = [
    ./configuration.nix
+   ./hardware-configuration.nix
];
```

Fix "file not found" error

```bash
$ git add .
```

And rebuild from the flake

```bash
$ pwd # make sure you are in the right directory
/home/yourname/dotfiles
$ sudo nixos-rebuild switch --flake .#yourhostname
```

# Home manager

Add home-manager input to your flake

```diff
  inputs = {
    nixpkgs = { url = "github:nixos/nixpkgs/nixos-unstable"; };
+   home-manager = {
+     url = github:nix-community/home-manager;
+     inputs.nixpkgs.follows = "nixpkgs";
+   };
  };
```

and a NixOS module

```diff
    nixosConfigurations = {
      yourhostname = inputs.nixpkgs.lib.nixosSystem {
        system = "x86_64-linux";
        modules = [
          ./configuration.nix
          ./hardware-configuration.nix
+         inputs.home-manager.nixosModules.home-manager
+         {
+           home-manager = {
+             useUserPackages = true;
+             useGlobalPkgs = true;
+             extraSpecialArgs = {inherit inputs;};
+             users.yourname = ./home.nix;
+           };
+         }
        ];
        specialArgs = { inherit inputs; };
```

Now we need to create a _home.nix_ file we just called. And for this example, we
will configure zsh

```nix
# remember to git add home.nix!
{pkgs, ...}: {
  home.stateVersion = "22.11";
  programs.zsh.enable = true;
}
```

Let's see if everything works by running `nix flake check`, it should return no errors.
Everything works? Great, now we can _rice_

```nix
{pkgs, ...}: {
  home.stateVersion = "22.11";

  # install some packages
  home.packages = with pkgs; [ firefox du-dust neovim ];

  programs.zsh = {
    enable = true;
    # maybe some shell aliases
    shellAliases = {
      vim = "nvim";
      du = "dust";
    };
    # completion would be nice
    enableCompletion = true;
    autosuggestions.enable = true;
  };
}
```

But, how do we know which hm modules are available? Just check in their docs
or [source code directly (it's more readable for some reason)](https://github.com/nix-community/home-manager/tree/master/modules/programs)

Even if you have no experience with nix, reading the source code should give you a good idea of what options are available and how to use them
[zsh example](https://github.com/nix-community/home-manager/blob/62cb5bcf93896e4dd6b4507dac7ba2e2e3abc9d7/modules/programs/zsh.nix#L213-L427)

Cool, but what if there is no module, like for example for neofetch?

We can just use _symlinks_

```nix
xdg.configFile."neofetch/config".source = ./neofetch-config-file;
# this makes a symlink to ~/.config/neofetch/config
# same as
home.file.".config/neofetch/config".source = ./neofetch-config-file;

# we can also use it with dirs, like this
xdg.configFile."nvim".source = ./nvim;

```

# Whats next?

Now you should have a basic understanding how NixOS works, at least enough to get you started.
My advise is: [sourcegraph is your friend](https://sourcegraph.com/search?q=context%3Aglobal+lang%3Anix+programs.neovim&patternType=standard&sm=1).
Also try reading other people's configurations, but without blidnly copy-pasting, thats not the point.

Here are some good config repos:

- [fufexan](https://github.com/fufexan/dotfiles)
- [viperML](https://github.com/viperML/dotfiles)
- [NobbZ](https://github.com/NobbZ/nixos-config)
- [gytis](https://github.com/gytis-ivaskevicius/nixfiles)
- [hlissner](https://github.com/hlissner/dotfiles)
- [notusknot](https://github.com/notusknot/dotfiles-nix)

[And heres mine <3](https://dotfiles.sioodmy.dev)

# Donate

I've spent way more time on this than I expected, so please
consider [donating](/donate) <3

Also if something is unclear/outdated/incorrect/there-is-a-better-way make sure to [reach me out](/about) and I'll try my best to keep it updated.
