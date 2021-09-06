<p align="center">
  
  <img src="https://media.giphy.com/media/c5Qr4yg3DRams/giphy.gif">
</p>

<h3 align="center">
  Booster
</h3>

<p align="center">
  The Buildpack-powered build service
</p>

---

# 🚀 Booster

This service handles image builds for the Magic Container Registry. It currently uses [kpack](https://github.com/pivotal/kpack) to handle the actual build process, but the plan is to drop that dependency in the future.

## Releasing

Before Spaceship is self-hosting, you can build an image manually with `pack`:

```sh
pack build us.gcr.io/onspaceship/booster --builder paketobuildpacks/builder:tiny
```
