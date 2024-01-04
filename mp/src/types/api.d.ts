declare namespace API {
  interface User {
    id?: number;
    name?: string;
    email?: string;
    openid?: string;
    avatar?: string;
    role?: number;
    status?: number;
    meta?: Record<string, any>;
    hasPassword?: boolean;
    createdAt?: string;
    updatedAt?: string;
  }

  interface Article {
    id?: number;
    userId?: string;
    title?: string;
    content?: string;
    textContent?: string;
    preview?: string;
    viewCount?: number;
    likeCount?: number;
    commentCount?: number;
    favoriteCount?: number;
    createdAt?: string;
    updatedAt?: string;
    hasLiked?: boolean;
    hasFavorited?: boolean;
    user?: User;
    tags?: Tag[];
    likes?: Like[];
    highlights?: Record<string, any[]>
  }

  interface Tag {
    id?: number;
    name?: string;
    sort?: number;
    createdAt?: string;
    updatedAt?: string;
  }

  interface Notification {
    id?: number;
    userId?: number;
    fromUserId?: number;
    type?: string;
    subject?: string;
    message?: string;
    data?: Record<string, any>;
    readAt?: string;
    createdAt?: string;
    updatedAt?: string;
    user?: User;
    fromUser?: User;
  }

  interface Like {
    id?: number;
    userId?: number;
    articleId?: number;
    createdAt?: string;
    updatedAt?: string;
    deletedAt?: string;
    article?: Article;
  }

  interface Comment {
    id?: number;
    userId?: number;
    articleId?: number;
    parentId?: number;
    commentId?: number;
    content?: string;
    textContent?: string;
    upvoteCount?: number;
    replyCount?: number;
    createdAt?: string;
    updatedAt?: string;
    deletedAt?: string;
    hasUpvoted?: boolean;
    user?: User;
    article?: Article;
    replies?: Comment[];
    parent?: Comment;
    Upvotes?: Upvote[];
  }

  interface Upvote {
    id?: number;
    userId?: number;
    commentId?: number;
    createdAt?: string;
    updatedAt?: string;
    deletedAt?: string;
    user?: User;
    comment?: Comment;
  }

  interface Favorite {
    id?: number;
    userId?: number;
    articleId?: number;
    createdAt?: string;
    updatedAt?: string;
    deletedAt?: string;
    user?: User;
    article?: Article;
  }
}
