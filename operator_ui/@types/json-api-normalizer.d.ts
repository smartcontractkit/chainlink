/* eslint @typescript-eslint/class-name-casing: 'off' */

declare module 'json-api-normalizer' {
  export interface JsonApiResponse<
    TData extends
      | ResourceObject<any, any, any, any>[]
      | ResourceObject<any, any, any, any>
      | null =
      | ResourceObject<any, any, any, any>[]
      | ResourceObject<any, any, any, any>,
    TError extends ErrorsObject[] = ErrorsObject[],
    TIncluded extends ResourceObject<any, any, any, any>[] | never = never,
    TMeta extends Record<string, any> | never = never,
    TLinks extends LinksObject | never = never,
  > {
    data: TData
    included?: TIncluded
    links: TLinks
    errors?: TError
    meta: TMeta
  }

  /**
   * Where specified, a links member can be used to represent links. The value of each links
   * member MUST be an object (a “links object”).
   *
   *   Each member of a links object is a “link”. A link MUST be represented as either:
   *
   *    - a string containing the link’s URL.
   *    - an object (“link object”) which can contain the following members:
   *    - href: a string containing the link’s URL.
   *    - meta: a meta object containing non-standard meta-information about the link.
   */
  type LinksObject = Partial<Record<string, string | LinkObject>>
  export interface LinkObject<
    TLinkMeta extends Record<string, any> = Record<string, any>,
  > {
    href: string
    meta: TLinkMeta
  }

  export interface ErrorsObject<
    TMeta extends Record<string, any> = Record<string, any>,
    TLinks extends LinksObject = LinksObject,
  > {
    id?: string
    links?: TLinks
    status?: string
    code?: string
    title?: string
    detail?: string
    source?: {
      pointer?: string
      parameter?: string
      [propName: string]: any
    }
    meta?: TMeta
  }

  /**
   * Resource Identifier Objects
   * A “resource identifier object” is an object that identifies an individual resource.
   *
   * A “resource identifier object” MUST contain type and id members.
   *
   * A “resource identifier object” MAY also include a meta member, whose value is a meta object
   * that contains non-standard meta-information.
   */
  export interface ResourceIdentifierObject<
    TMeta extends Record<string, any> = Record<string, any>,
  > {
    type: string
    id: string
    meta?: TMeta
  }

  /**
   * Resource Linkage
   *
   * Resource linkage in a compound document allows a client to
   *  link together all of the included resource objects without
   *  having to GET any URLs via links.
   *
   * Resource linkage MUST be represented as one of the following:
   *
   * - null for empty to-one relationships.
   * - an empty array ([]) for empty to-many relationships.
   * - a single resource identifier object for non-empty to-one relationships.
   * - an array of resource identifier objects for non-empty to-many relationships.
   */
  export type ResourceLinkage =
    | null
    | []
    | ResourceIdentifierObject
    | ResourceIdentifierObject[]

  /**
   * The value of the relationships key MUST be an object (a “relationships object”). Members of
   * the relationships object (“relationships”) represent references from the resource object in
   * which it’s defined to other resource objects.
   *
   * Relationships may be to-one or to-many.
   *
   * A “relationship object” MUST contain at least one of the following:
   *
   * - links: a links object containing at least one of the following:
   *    - self: a link for the relationship itself (a “relationship link”).
   *        This link allows the client to directly manipulate the relationship.
   *        For example, removing an author through an article’s relationship URL would
   *        disconnect the person from the article without deleting the people resource itself.
   *        When fetched successfully, this link returns the linkage for the related resources
   *        as its primary data. (See Fetching Relationships.)
   *    - related: a related resource link
   *- data: resource linkage
   *- meta: a meta object that contains non-standard meta-information about the relationship.
   *
   * A relationship object that represents a to-many relationship MAY also contain pagination
   * links under the links member, as described below.
   * Any pagination links in a relationship object MUST paginate the relationship data,
   * not the related resources.
   */
  export type Relationship<
    TMeta extends Record<string, any> = Record<string, any>,
    TLinks extends LinksObject = LinksObject,
  > = _Relationship<TMeta, TLinks>
  export interface _Relationship<
    TMeta extends Record<string, any>,
    TLinks extends LinksObject,
  > {
    data?: JsonApiResponse | JsonApiResponse[]
    links?: TLinks
    meta?: TMeta
  }

  /**
   * The value of the attributes key MUST be an object (an “attributes object”). Members of the
   * attributes object (“attributes”) represent information about the resource object in which
   * it’s defined.
   *
   * Attributes may contain any valid JSON value.
   *
   * Complex data structures involving JSON objects and arrays are allowed as attribute values.
   * However, any object that constitutes or is contained in an attribute MUST NOT contain a
   * relationships or links member, as those members are reserved by this specification for
   * future use.
   *
   * Although has-one foreign keys (e.g. author_id) are often stored internally alongside other
   * information to be represented in a resource object, these keys SHOULD NOT appear as
   * attributes.
   */
  interface AttributesObject extends Record<string, any> {}

  /**
   * Resource Objects
   *
   * “Resource objects” appear in a JSON:API document to represent resources.
   *
   *  A resource object MUST contain at least the following top-level members:
   * - id
   * - type
   *
   * Exception: The id member is not required when the resource object originates at the client
   * and represents a new resource to be created on the server.
   *
   *  In addition, a resource object MAY contain any of these top-level members:
   * - attributes: an attributes object representing some of the resource’s data.
   * - relationships: a relationships object describing relationships between the resource and other JSON:API resources.
   * - links: a links object containing links related to the resource.
   * - meta: a meta object containing non-standard meta-information about a resource that can not be represented as an attribute or relationship.
   */
  export interface ResourceObject<
    TAttributes extends AttributesObject | never = never,
    TRelationships extends
      | Record<string, Relationship<Record<string, any>>>
      | never = never,
    TMeta extends Record<string, any> | never = never,
    TLinks extends LinksObject | never = never,
  > {
    id: string
    type: string
    attributes: TAttributes
    relationships: TRelationships
    links: TLinks
    meta: TMeta
  }

  // we cant infer TNormalized from the arguments
  // because of the transformations it does on the supplied data
  // unfortunately, this means that the user will have to supply
  // their own typing for TNormalized

  /**
   * Normalize JSON:API spec compliant JSON
   *
   * @param json The JSON:API spec compliant JSON to normalize
   * @param opts Options for normalizing
   */
  export default function normalize<TNormalized>(
    json: JsonApiResponse<any, any, any, any, any>,
    opts?: Opts,
  ): TNormalized

  interface Opts {
    camelizeKeys?: boolean
    camelizeTypeValues?: boolean
    endpoint?: string
    filterEndpoint?: boolean
  }
}
