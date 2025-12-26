# ShiyuAdmin - Système d'administration générique (FR, version courte)

> Auteur : Wang Shiyu (JavaPub)  
> Site officiel : https://javapub.net.cn/  
> Dépôt : **ShiyuAdmin** (open-source sur GitHub, avec miroir Gitee pour un accès plus rapide en Chine)
> - GitHub (dépôt principal) : `https://github.com/Rodert/ShiyuAdmin`  
> - Gitee (miroir chinois) : `https://gitee.com/rodert/ShiyuAdmin`

---

## 1. Présentation

ShiyuAdmin est un modèle de back-office moderne basé sur :

- un backend Go (Gin + Gorm)
- un frontend React + Ant Design Pro

Il peut servir de :

- point de départ pour un panneau d'administration générique
- projet d'exemple pour apprendre Go + React + RBAC
- base pour des outils internes ou des projets personnels

Pour une documentation complète (en chinois), consultez `README.md` à la racine du dépôt.

---

## 2. Stack technique (résumé)

- Backend
  - Go 1.23+
  - Gin, Gorm
  - PostgreSQL / MySQL / SQLite
  - Redis
  - Authentification JWT + modèle de permissions RBAC

- Frontend
  - React 19
  - Umi Max
  - Ant Design & Ant Design Pro Components

---

## 3. Démarrage rapide (Docker recommandé)

Prérequis :

- Docker et Docker Compose installés

Étapes :

```bash
# Cloner le dépôt (remplacer par votre URL de dépôt si vous avez fait un fork)
git clone https://github.com/Rodert/ShiyuAdmin.git
cd ShiyuAdmin

docker-compose up -d
```

Services exposés :

- Backend : `http://localhost:8080`
- Frontend : `http://localhost:8000`
- Compte administrateur par défaut : `admin` / `Admin@123`

Les autres modes de lancement (développement local, intégration dans un projet existant, etc.)
sont décrits dans le `README.md` chinois.

---

## 4. Contributions & support

- **Mettez une étoile** : N'hésitez pas à mettre une étoile (★) au dépôt si le projet vous aide  
- **Signaler un bug** : [Créer une issue](https://github.com/Rodert/ShiyuAdmin/issues/new?labels=bug)  
- **Proposer une fonctionnalité** : [Créer une issue](https://github.com/Rodert/ShiyuAdmin/issues/new?labels=enhancement)  
- **Soumettre une PR** : Améliorations de code, nouvelles fonctionnalités, mises à jour de documentation  
- **Contact** : Suivez le compte WeChat « JavaPub » pour contacter l'auteur (en chinois)
