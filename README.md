# medivpn

medivpn is my personal VPN configuration tool. It's designed to simplify the act of switching VPN servers and toggling my VPN on/off for my particular setup, without having to ssh into relevant machines and intermittently deal with Tailscale's CLI.

As for the networking setup that necessitates this: I have a dedicated server running 24/7 offering an exit node that routes all non-local traffic through wireguard to commercial VPN servers. All any of my devices have to do to 'opt in' is select that dedicated server's exit node through Tailscale.

## Why...?

Android.

By routing my traffic through a single machine, it allows me to use both Tailscale and commercial VPN services on Android devices like my phone. Android otherwise limits users to a single VPN running at a time, meaning I would have to choose between using Tailscale (and my home lab) or using commercial VPN servers. This setup therefore gives me the best of both worlds while also serving to simplify configuration (WireGuard configs are only present on the dedicated server and can be switched between without ssh-ing in) and setup for future devices (all I have to do is install Tailscale and this CLI tool).

Since the commercial VPN servers I'm using rate limit me anyway, and are therefore not appropriate for large downloads with high volumes of network traffic, I've also adopted the same setup on both my desktop and my laptop.

## Project Structure

No. This repo isn't big enough for that ._.

## Limitations

The CLI is exclusively for my desktop and laptop at the moment. To switch servers on my laptop I'd have to ssh into an existing machine via Termux, which kind of defeats the purpose of making this CLI in the first place. I could 